package container

import (
	"context"
	"errors"
	"fmt"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker/components/image"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/pathutil"
	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/blkiodev"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/joho/godotenv"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
)

var _ core.Lifecycle = (*Container)(nil)

func (c *Container) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	var spec Spec
	if err := c.LoadSpec(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}
	seen := make(map[string]bool)
	deps := make([]string, 0, len(spec.DependsOn.Items)+len(spec.Links.Values()))
	addDep := func(service string) {
		if seen[service] {
			return
		}
		seen[service] = true
		deps = append(deps, service)
	}
	for _, dep := range spec.DependsOn.Items {
		addDep(dep.Service.Value)
	}
	for _, link := range spec.Links {
		addDep(link.Service)
	}
	return &core.DependenciesOutput{Components: deps}, nil
}

func (c *Container) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	var spec Spec
	if err := c.LoadSpec(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	// NOTE [IMAGE_SUBCOMPONENT]: Should create image as subcomponent instead of
	// copying spec in to state.
	c.State.Image.Spec = yamlutil.MustMarshalString(image.Spec{
		Platform: spec.Platform.Value,
		Build:    spec.Build,
	})

	if err := c.ensureImage(ctx, &spec); err != nil {
		return nil, fmt.Errorf("ensuring image: %w", err)
	}

	if err := c.removeExistingContainerByName(ctx, spec.ContainerName.Value); err != nil {
		return nil, fmt.Errorf("removing existing container %q: %w", spec.ContainerName, err)
	}

	if err := c.create(ctx, &spec); err != nil {
		return nil, fmt.Errorf("creating container: %w", err)
	}

	if err := c.start(ctx); err != nil {
		c.Logger.Infof("starting container %q: %v", c.State.ContainerID, err)
	}

	return &core.InitializeOutput{}, nil
}

func (c *Container) create(ctx context.Context, spec *Spec) error {
	var healthCfg *container.HealthConfig
	if spec.Healthcheck != nil {
		healthCfg = &container.HealthConfig{
			Test:        strslice.StrSlice(spec.Healthcheck.Test.Parts.Values()),
			Interval:    spec.Healthcheck.Interval.Duration,
			Timeout:     spec.Healthcheck.Timeout.Duration,
			Retries:     spec.Healthcheck.Retries.Int(),
			StartPeriod: spec.Healthcheck.StartPeriod.Duration,
		}
	}

	labels := spec.Labels.Map()
	for k, v := range c.GetExoLabels() {
		labels[k] = v
	}

	envMap := map[string]string{}
	for _, envFilePath := range spec.EnvFile.Items {
		if !path.IsAbs(envFilePath.Value) {
			envFilePath.Value = path.Join(c.WorkspaceRoot, envFilePath.Value)
		}
		if !pathutil.HasPathPrefix(envFilePath.Value, c.WorkspaceRoot) {
			return fmt.Errorf("env file %s is not contained within the workspace", envFilePath.Value)
		}
		envFileVars, err := godotenv.Read(envFilePath.Value)
		if err != nil {
			return fmt.Errorf("reading env file %s: %w", envFilePath.Value, err)
		}
		for k, v := range envFileVars {
			envMap[k] = v
		}
	}
	for _, item := range spec.Environment.Items {
		if item.Value == "" {
			if v, ok := c.WorkspaceEnvironment[item.Key]; ok {
				envMap[item.Key] = v
			}
		} else {
			envMap[item.Key] = item.Value
		}
	}
	envSlice := []string{}
	for k, v := range envMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	containerCfg := &container.Config{
		Hostname:     spec.Hostname.Value,
		Domainname:   spec.Domainname.Value,
		User:         spec.User.Value,
		ExposedPorts: make(nat.PortSet),
		Tty:          spec.TTY.Value,
		OpenStdin:    spec.StdinOpen.Value,
		// StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
		Env:         envSlice,
		Healthcheck: healthCfg,
		// ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (meaning treat as a command line) (Windows specific).

		Image: c.State.Image.ID,
		// Volumes         map[string]struct{} // List of volumes (mounts) used for the container
		WorkingDir: spec.WorkingDir.Value,
		Entrypoint: strslice.StrSlice(spec.Entrypoint.Parts.Values()),
		// NetworkDisabled bool                `json:",omitempty"` // Is network disabled
		MacAddress: spec.MacAddress.Value,
		// OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
		Labels:     labels,
		StopSignal: spec.StopSignal.Value,
		// Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}

	if spec.Command.IsShellForm {
		containerCfg.Cmd = append(append([]string{}, c.State.Image.Shell...), spec.Command.Parts[0].Value)
	} else {
		containerCfg.Cmd = spec.Command.Parts.Values()
	}
	if len(containerCfg.Cmd) == 0 {
		containerCfg.Cmd = c.State.Image.Command
	}

	if len(containerCfg.Entrypoint) == 0 {
		containerCfg.Entrypoint = c.State.Image.Entrypoint
	}

	if containerCfg.WorkingDir == "" {
		containerCfg.WorkingDir = c.State.Image.WorkingDir
	}

	if spec.StopGracePeriod != nil {
		timeout := int(spec.StopGracePeriod.Duration.Round(time.Second).Seconds())
		containerCfg.StopTimeout = &timeout
	}

	exposePort := func(min, max uint16, protocol string) {
		for n := min; n <= max; n++ {
			port := nat.Port(compose.FormatPort(n, protocol))
			containerCfg.ExposedPorts[port] = struct{}{}
		}
	}
	for _, exposed := range spec.Expose {
		exposePort(exposed.Min, exposed.Max, exposed.Protocol)
	}
	for _, mapping := range spec.Ports {
		exposePort(mapping.Target.Min, mapping.Target.Max, mapping.Protocol)
	}

	logCfg := container.LogConfig{}
	if spec.Logging.Driver.Value == "" && len(spec.Logging.Options.Items) == 0 {
		// No logging configuration specified, so default to logging to exo's
		// syslog service.
		logCfg.Type = "syslog"
		bridge, err := c.Docker.NetworkInspect(ctx, "bridge", types.NetworkInspectOptions{})
		if err != nil {
			return fmt.Errorf("inspecting bridge network: %w", err)
		}
		if len(bridge.IPAM.Config) != 1 {
			return fmt.Errorf("bridge network has %d IPAM configs, expected 1", len(bridge.IPAM.Config))
		}
		syslogHost := bridge.IPAM.Config[0].Gateway

		logCfg.Config = map[string]string{
			"syslog-address":  fmt.Sprintf("udp://%s:%d", syslogHost, c.SyslogPort),
			"syslog-facility": "1", // "user-level messages"
			"tag":             c.ComponentID,
			"syslog-format":   "rfc5424micro",
		}
	} else {
		logCfg.Type = spec.Logging.Driver.Value
		logCfg.Config = spec.Logging.Options.Map()
	}

	blkioWeightDevice := make([]*blkiodev.WeightDevice, len(spec.BlkioConfig.WeightDevice))
	for i, weightDevice := range spec.BlkioConfig.WeightDevice {
		blkioWeightDevice[i] = &blkiodev.WeightDevice{
			Path:   weightDevice.Path.Value,
			Weight: weightDevice.Weight.Uint16(),
		}
	}

	hostCfg := &container.HostConfig{
		//// Applicable to all platforms
		//Binds           []string      // List of volume bindings for this container
		//ContainerIDFile string        // File (path) where the containerId is written
		LogConfig: logCfg,
		//NetworkMode     NetworkMode   // Network mode to use for the container
		PortBindings: make(nat.PortMap),
		// TODO: Potentially inherit from deploy's restart_policy.
		RestartPolicy: container.RestartPolicy{
			Name: spec.Restart.Value,
		},
		//AutoRemove      bool          // Automatically remove container when it exits
		//VolumeDriver    string        // Name of the volume driver used to mount volumes

		//// Applicable to UNIX platforms
		CapAdd:  spec.CapAdd.Values(),
		CapDrop: spec.CapDrop.Values(),
		//CgroupnsMode    CgroupnsMode      // Cgroup namespace mode to use for the container
		DNS:        spec.DNS.Values(),
		DNSOptions: spec.DNSOptions.Values(),
		DNSSearch:  spec.DNSSearch.Values(),
		ExtraHosts: spec.ExtraHosts.Values(),
		GroupAdd:   spec.GroupAdd.Values(),
		//Cgroup          CgroupSpec        // Cgroup to use for the container

		// See NOTE: [RESOLVING SERVICE CONTAINERS].
		Links: append(append([]string{}, spec.Links.Values()...), spec.ExternalLinks.Values()...),
		//OomScoreAdj     int               // Container preference for OOM-killing
		Privileged: spec.Privileged.Value,
		//PublishAllPorts bool              // Should docker publish all exposed port for the container
		ReadonlyRootfs: spec.ReadOnly.Value,
		SecurityOpt:    spec.SecurityOpt.Values(),
		StorageOpt:     spec.StorageOpt.Map(),
		//UTSMode         UTSMode           // UTS namespace to use for the container
		UsernsMode: container.UsernsMode(spec.UsernsMode.Value),
		ShmSize:    spec.ShmSize.Int64(),
		Sysctls:    spec.Sysctls.Map(),
		Runtime:    spec.Runtime.Value,

		//// Applicable to Windows
		//ConsoleSize [2]uint   // Initial console size (height,width)

		//// Contains container's resources (cgroups, ulimits)
		Resources: container.Resources{
			CPUCount:             spec.CPUCount.Value,
			CPUPercent:           spec.CPUPercent.Value,
			CPUShares:            spec.CPUShares.Value,
			CPUPeriod:            spec.CPUPeriod.Value,
			CPUQuota:             spec.CPUQuota.Value,
			Memory:               spec.MemoryLimit.Int64(),
			MemoryReservation:    spec.MemoryReservation.Int64(),
			MemorySwappiness:     spec.MemorySwappiness.Int64Ptr(),
			MemorySwap:           spec.MemswapLimit.Int64(),
			CPURealtimePeriod:    spec.CPURealtimePeriod.Duration.Microseconds(),
			CPURealtimeRuntime:   spec.CPURealtimeRuntime.Duration.Microseconds(),
			BlkioWeight:          spec.BlkioConfig.Weight.Uint16(),
			BlkioWeightDevice:    blkioWeightDevice,
			BlkioDeviceReadBps:   convertThrottleDevice(spec.BlkioConfig.DeviceReadBPS),
			BlkioDeviceReadIOps:  convertThrottleDevice(spec.BlkioConfig.DeviceReadIOPS),
			BlkioDeviceWriteBps:  convertThrottleDevice(spec.BlkioConfig.DeviceWriteBPS),
			BlkioDeviceWriteIOps: convertThrottleDevice(spec.BlkioConfig.DeviceWriteIOPS),
			CpusetCpus:           spec.CPUSet.Value,
			CgroupParent:         spec.CgroupParent.Value,
			DeviceCgroupRules:    spec.DeviceCgroupRules.Values(),
			Devices:              convertDeviceMappings(spec.Devices),
			OomKillDisable:       spec.OomKillDisable.Ptr(),
			PidsLimit:            spec.PidsLimit.Int64Ptr(),
			Ulimits:              convertUlimits(spec.Ulimits),
		},

		OomScoreAdj: spec.OomScoreAdj.Int(),

		//// MaskedPaths is the list of paths to be masked inside the container (this overrides the default set of paths)
		//MaskedPaths []string

		//// ReadonlyPaths is the list of paths to be set as read-only inside the container (this overrides the default set of paths)
		//ReadonlyPaths []string

		//// Run a custom init inside the container, if null, use the daemon's configured settings
		Init: spec.Init.Ptr(),
	}

	if hostCfg.IpcMode, err = c.parseIPCMode(spec.IPC.Value); err != nil {
		return err
	}

	if hostCfg.Isolation, err = parseIsolation(spec.Isolation.Value); err != nil {
		return err
	}

	if hostCfg.NetworkMode, err = c.parseNetworkMode(spec.NetworkMode.Value); err != nil {
		return err
	}

	if hostCfg.PidMode, err = c.parsePIDMode(spec.PidMode.Value); err != nil {
		return err
	}

	if hostCfg.VolumesFrom, err = c.parseVolumesFrom(spec.VolumesFrom.Values()); err != nil {
		return err
	}

	if len(spec.Tmpfs.Items) > 0 {
		hostCfg.Tmpfs = make(map[string]string, len(spec.Tmpfs.Items))
		// This matches the docker-compose behaviour for specifying tmpfs mounts with the service-level `tmpfs` option.
		for _, path := range spec.Tmpfs.Items {
			hostCfg.Tmpfs[path.Value] = ""
		}
	}

	// TODO: make the user home directory a parameter of the container.
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not get user %w", err)
	}
	userHomeDir := user.HomeDir

	hostCfg.Mounts = make([]mount.Mount, len(spec.Volumes))
	for i, v := range spec.Volumes {
		mnt, err := makeMountFromVolumeMount(c.WorkspaceRoot, userHomeDir, v)
		if err != nil {
			return fmt.Errorf("invalid mount at index %d: %w", i, err)
		}
		hostCfg.Mounts[i] = mnt
	}

	for _, mapping := range spec.Ports {
		targetLow, targetHigh := int(mapping.Target.Min), int(mapping.Target.Max)
		for targetPort := targetLow; targetPort <= targetHigh; targetPort += 1 {
			publishedLow, publishedHigh := int(mapping.Published.Min), int(mapping.Published.Max)
			publishedDiff, targetDiff := publishedHigh-publishedLow, targetHigh-targetLow
			if publishedDiff != 1 && publishedDiff != targetDiff {
				return fmt.Errorf("unexpected number of ports")
			}

			target := nat.Port(strconv.Itoa(targetPort))
			hostPort := strconv.Itoa(publishedLow + publishedDiff)
			if mapping.Protocol != "" {
				hostPort += "/" + mapping.Protocol
			}
			bindings := hostCfg.PortBindings[target]
			bindings = append(bindings, nat.PortBinding{
				HostIP:   mapping.HostIP,
				HostPort: hostPort,
			})

			// TODO: Handle mapping.Mode
			hostCfg.PortBindings[target] = bindings
		}
	}

	networkCfg := &network.NetworkingConfig{
		EndpointsConfig: make(map[string]*network.EndpointSettings), // Endpoint configs for each connecting network
	}
	// Docker only allows a single network to be specified when creating a container. The other networks must be
	// connected after the container is started. See https://github.com/moby/moby/issues/29265#issuecomment-265909198.
	var remainingNetworks []compose.ServiceNetwork
	if len(spec.Networks.Items) > 0 {
		firstNetwork := spec.Networks.Items[0]
		remainingNetworks = spec.Networks.Items[1:]
		networkCfg.EndpointsConfig[firstNetwork.Key] = c.endpointSettings(firstNetwork, spec)
	}

	var platform *v1.Platform
	//platform := &v1.Platform{
	//	//// Architecture field specifies the CPU architecture, for example
	//	//// `amd64` or `ppc64`.
	//	//Architecture string `json:"architecture"`

	//	//// OS specifies the operating system, for example `linux` or `windows`.
	//	//OS string `json:"os"`

	//	//// OSVersion is an optional field specifying the operating system
	//	//// version, for example on Windows `10.0.14393.1066`.
	//	//OSVersion string `json:"os.version,omitempty"`

	//	//// OSFeatures is an optional field specifying an array of strings,
	//	//// each listing a required OS feature (for example on Windows `win32k`).
	//	//OSFeatures []string `json:"os.features,omitempty"`

	//	//// Variant is an optional field specifying a variant of the CPU, for
	//	//// example `v7` to specify ARMv7 when architecture is `arm`.
	//	//Variant string `json:"variant,omitempty"`
	//}
	createdBody, err := c.Docker.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, platform, spec.ContainerName.Value)
	if err != nil {
		return err
	}
	c.State.ContainerID = createdBody.ID
	var netConnects errgroup.Group
	for _, network := range remainingNetworks {
		network := network
		netConnects.Go(func() error {
			return c.Docker.NetworkConnect(ctx, network.Key, createdBody.ID, c.endpointSettings(network, spec))
		})
	}

	return netConnects.Wait()
}

func (c *Container) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	{
		// NOTE [MIGRATE_CONTAINER_STATE]: Data migration that copies additional
		// information from spec in to state.  Can be removed after sufficient time
		// has passed from October 2021, or when the [IMAGE_SUBCOMPONENT] note has
		// been resolved.
		var spec Spec
		if err := c.LoadSpec(input.Spec, &spec); err != nil {
			return nil, fmt.Errorf("loading spec: %w", err)
		}
		c.State.Image.Spec = yamlutil.MustMarshalString(image.Spec{
			Platform: spec.Platform.Value,
			Build:    spec.Build,
		})
	}

	if c.State.ContainerID == "" {
		c.State.Running = false
	} else {
		inspection, err := c.Docker.ContainerInspect(ctx, c.State.ContainerID)
		if err != nil {
			return nil, fmt.Errorf("inspecting container: %w", err)
		}

		c.State.Running = inspection.State.Running
	}
	return &core.RefreshOutput{}, nil
}

func (c *Container) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	if c.State.ContainerID == "" {
		return &core.DisposeOutput{}, nil
	}

	if err := c.stop(ctx, nil); err != nil {
		c.Logger.Infof("stopping container %q: %v", c.State.ContainerID, err)
	}
	err := c.Docker.ContainerRemove(ctx, c.State.ContainerID, types.ContainerRemoveOptions{
		// XXX RemoveVolumes: ???,
		// XXX RemoveLinks: ???,
		Force: true, // OK?
	})
	if docker.IsErrNotFound(err) {
		c.Logger.Infof("container to be removed not found: %q", c.State.ContainerID)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	c.State.ContainerID = ""
	return &core.DisposeOutput{}, nil
}

func (c *Container) removeExistingContainerByName(ctx context.Context, name string) error {
	// If a container with this name already exists, remove it.
	containers, err := c.Docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: name,
		}),
		All: true,
	})
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}
	if len(containers) > 0 {
		if err := c.Docker.ContainerRemove(ctx, containers[0].ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) endpointSettings(sn compose.ServiceNetwork, spec *Spec) *network.EndpointSettings {
	es := &network.EndpointSettings{
		Aliases: append([]string{c.ComponentID, c.ComponentName}, sn.Aliases.Values()...),
		Links:   append(append([]string{}, spec.Links.Values()...), spec.ExternalLinks.Values()...),
	}

	if sn.IPV4Address.Value != "" || sn.IPV6Address.Value != "" || len(sn.LinkLocalIPs) > 0 {
		es.IPAMConfig = &network.EndpointIPAMConfig{
			IPv4Address:  sn.IPV4Address.Value,
			IPv6Address:  sn.IPV6Address.Value,
			LinkLocalIPs: sn.LinkLocalIPs.Values(),
		}
	}

	return es
}

func convertThrottleDevice(in []compose.ThrottleDevice) []*blkiodev.ThrottleDevice {
	if in == nil {
		return nil
	}

	out := make([]*blkiodev.ThrottleDevice, len(in))
	for i, throttleDevice := range in {
		out[i] = &blkiodev.ThrottleDevice{
			Path: throttleDevice.Path.Value,
			Rate: throttleDevice.Rate.Uint64(),
		}
	}
	return out
}

func convertDeviceMappings(in []compose.DeviceMapping) []container.DeviceMapping {
	if in == nil {
		return nil
	}

	out := make([]container.DeviceMapping, len(in))
	for i, deviceMapping := range in {
		out[i] = container.DeviceMapping{
			PathOnHost:        deviceMapping.PathOnHost,
			PathInContainer:   deviceMapping.PathInContainer,
			CgroupPermissions: deviceMapping.CgroupPermissions,
		}
	}
	return out
}

func convertUlimits(in compose.Ulimits) []*units.Ulimit {
	if in == nil {
		return nil
	}

	out := make([]*units.Ulimit, len(in))
	for i, ulimit := range in {
		out[i] = &units.Ulimit{
			Name: ulimit.Name,
			Hard: ulimit.Hard.Value,
			Soft: ulimit.Soft.Value,
		}
	}
	return out
}

func (c *Container) parseIPCMode(in string) (container.IpcMode, error) {
	switch in {
	case "", "none", "private", "shareable", "host":
		return container.IpcMode(in), nil
	}

	if strings.HasPrefix(in, "container:") {
		return container.IpcMode(in), nil
	}

	// Note that all services that share an IPC namespace with another service actually shares a namespace
	// with that service's first container. See https://github.com/docker/compose/blob/v2.0.0-rc.3/compose/service.py#L1379.
	if strings.HasPrefix(in, "service:") {
		// XXX: We need to be able to look up the container name(s) for the referenced service to be able to
		// convert this to the `container:name` format supported by Docker.
		return container.IpcMode(""), errors.New("service IPC mode not yet supported")
	}

	return container.IpcMode(""), fmt.Errorf("unsupported IPC mode: %q", in)
}

func parseIsolation(in string) (container.Isolation, error) {
	switch in {
	case "", "default", "process", "hyperv":
		return container.Isolation(in), nil

	default:
		return container.IsolationEmpty, fmt.Errorf("invalid isolation: %q", in)
	}
}

func (c *Container) parseNetworkMode(in string) (container.NetworkMode, error) {
	if strings.HasPrefix(in, "service:") {
		return container.NetworkMode(""), errors.New("service network mode not yet supported")
	}
	return container.NetworkMode(in), nil
}

func (c *Container) parsePIDMode(in string) (container.PidMode, error) {
	switch in {
	case "", "host":
		return container.PidMode(in), nil
	}

	if strings.HasPrefix(in, "container:") {
		return container.PidMode(in), nil
	}

	return container.PidMode(""), fmt.Errorf("Invalid PID mode: %q", in)
}

func (c *Container) parseVolumesFrom(in []string) ([]string, error) {
	for _, v := range in {
		if !strings.HasPrefix(v, "container:") {
			return nil, errors.New("volumes from service not yet supported")
		}
	}

	return in, nil
}
