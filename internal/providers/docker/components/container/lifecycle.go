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

func (c *Container) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	var spec Spec
	if err := yamlutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	// NOTE [IMAGE_SUBCOMPONENT]: Should create image as subcomponent instead of
	// copying spec in to state.
	c.State.Image.Spec = yamlutil.MustMarshalString(image.Spec{
		Platform: spec.Platform,
		Build:    spec.Build,
	})

	if err := c.ensureImage(ctx, &spec); err != nil {
		return nil, fmt.Errorf("ensuring image: %w", err)
	}

	if err := c.removeExistingContainerByName(ctx, spec.ContainerName); err != nil {
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
			Test:        strslice.StrSlice(spec.Healthcheck.Test.Parts),
			Interval:    time.Duration(spec.Healthcheck.Interval),
			Timeout:     time.Duration(spec.Healthcheck.Timeout),
			Retries:     spec.Healthcheck.Retries,
			StartPeriod: time.Duration(spec.Healthcheck.StartPeriod),
		}
	}

	labels := spec.Labels.WithoutNils()
	for k, v := range c.GetExoLabels() {
		labels[k] = v
	}

	envMap := map[string]string{}
	for _, envFilePath := range spec.EnvFile {
		if !path.IsAbs(envFilePath) {
			envFilePath = path.Join(c.WorkspaceRoot, envFilePath)
		}
		if !pathutil.HasPathPrefix(envFilePath, c.WorkspaceRoot) {
			return fmt.Errorf("env file %s is not contained within the workspace", envFilePath)
		}
		envFileVars, err := godotenv.Read(envFilePath)
		if err != nil {
			return fmt.Errorf("reading env file %s: %w", envFilePath, err)
		}
		for k, v := range envFileVars {
			envMap[k] = v
		}
	}
	for k, v := range spec.Environment {
		if v == nil {
			if v, ok := c.WorkspaceEnvironment[k]; ok {
				envMap[k] = v
			}
		} else {
			envMap[k] = *v
		}
	}
	envSlice := []string{}
	for k, v := range envMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	containerCfg := &container.Config{
		Hostname:     spec.Hostname,
		Domainname:   spec.Domainname,
		User:         spec.User,
		ExposedPorts: make(nat.PortSet),
		Tty:          spec.TTY,
		OpenStdin:    spec.StdinOpen,
		// StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
		Env:         envSlice,
		Healthcheck: healthCfg,
		// ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (meaning treat as a command line) (Windows specific).

		Image: c.State.Image.ID,
		// Volumes         map[string]struct{} // List of volumes (mounts) used for the container
		WorkingDir: spec.WorkingDir,
		Entrypoint: strslice.StrSlice(spec.Entrypoint.Parts),
		// NetworkDisabled bool                `json:",omitempty"` // Is network disabled
		MacAddress: spec.MacAddress,
		// OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
		Labels:     labels,
		StopSignal: spec.StopSignal,
		// Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}

	if spec.Command.IsShellForm {
		containerCfg.Cmd = append(append([]string{}, c.State.Image.Shell...), spec.Command.Parts[0])
	} else {
		containerCfg.Cmd = spec.Command.Parts
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
		timeout := int(time.Duration(*spec.StopGracePeriod).Round(time.Second).Seconds())
		containerCfg.StopTimeout = &timeout
	}

	exposePort := func(numbers string, protocol string) error {
		rng, err := compose.ParsePortRange(numbers, protocol)
		if err != nil {
			return fmt.Errorf("parsing port: %w", err)
		}
		for n := rng.Min; n <= rng.Max; n++ {
			port := nat.Port(compose.FormatPort(n, rng.Protocol))
			containerCfg.ExposedPorts[port] = struct{}{}
		}
		return nil
	}

	for _, exposed := range spec.Expose {
		if err := exposePort(exposed.Target, exposed.Protocol); err != nil {
			return fmt.Errorf("exposing port %q: %w", exposed.Target, err)
		}
	}
	for _, mapping := range spec.Ports {
		if err := exposePort(mapping.Target, mapping.Protocol); err != nil {
			return fmt.Errorf("exposing mapped port %q: %w", mapping.Target, err)
		}
	}

	logCfg := container.LogConfig{}
	if spec.Logging.Driver == "" && (spec.Logging.Options == nil || len(spec.Logging.Options) == 0) {
		// No logging configuration specified, so default to logging to exo's
		// syslog service.
		logCfg.Type = "syslog"
		logCfg.Config = map[string]string{
			"syslog-address":  fmt.Sprintf("udp://localhost:%d", c.SyslogPort),
			"syslog-facility": "1", // "user-level messages"
			"tag":             c.ComponentID,
			"syslog-format":   "rfc5424micro",
		}
	} else {
		logCfg.Type = spec.Logging.Driver
		logCfg.Config = spec.Logging.Options
	}

	blkioWeightDevice := make([]*blkiodev.WeightDevice, len(spec.BlkioConfig.WeightDevice))
	for i, weightDevice := range spec.BlkioConfig.WeightDevice {
		blkioWeightDevice[i] = &blkiodev.WeightDevice{
			Path:   weightDevice.Path,
			Weight: weightDevice.Weight,
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
			Name: spec.Restart,
		},
		//AutoRemove      bool          // Automatically remove container when it exits
		//VolumeDriver    string        // Name of the volume driver used to mount volumes

		//// Applicable to UNIX platforms
		CapAdd:  spec.CapAdd,
		CapDrop: spec.CapDrop,
		//CgroupnsMode    CgroupnsMode      // Cgroup namespace mode to use for the container
		DNS:        spec.DNS,
		DNSOptions: spec.DNSOptions,
		DNSSearch:  spec.DNSSearch,
		ExtraHosts: spec.ExtraHosts,
		GroupAdd:   spec.GroupAdd,
		//Cgroup          CgroupSpec        // Cgroup to use for the container

		// See NOTE: [RESOLVING SERVICE CONTAINERS].
		Links: append(append([]string{}, spec.Links...), spec.ExternalLinks...),
		//OomScoreAdj     int               // Container preference for OOM-killing
		Privileged: spec.Privileged,
		//PublishAllPorts bool              // Should docker publish all exposed port for the container
		ReadonlyRootfs: spec.ReadOnly,
		SecurityOpt:    spec.SecurityOpt,
		StorageOpt:     spec.StorageOpt,
		//UTSMode         UTSMode           // UTS namespace to use for the container
		UsernsMode: container.UsernsMode(spec.UsernsMode),
		ShmSize:    int64(spec.ShmSize),
		Sysctls:    spec.Sysctls.WithoutNils(),
		Runtime:    spec.Runtime,

		//// Applicable to Windows
		//ConsoleSize [2]uint   // Initial console size (height,width)

		//// Contains container's resources (cgroups, ulimits)
		Resources: container.Resources{
			CPUCount:             spec.CPUCount,
			CPUPercent:           spec.CPUPercent,
			CPUShares:            spec.CPUShares,
			CPUPeriod:            spec.CPUPeriod,
			CPUQuota:             spec.CPUQuota,
			Memory:               int64(spec.MemoryLimit),
			MemoryReservation:    int64(spec.MemoryReservation),
			MemorySwappiness:     spec.MemorySwappiness,
			MemorySwap:           int64(spec.MemswapLimit),
			CPURealtimePeriod:    time.Duration(spec.CPURealtimePeriod).Microseconds(),
			CPURealtimeRuntime:   time.Duration(spec.CPURealtimeRuntime).Microseconds(),
			BlkioWeight:          uint16(spec.BlkioConfig.Weight),
			BlkioWeightDevice:    blkioWeightDevice,
			BlkioDeviceReadBps:   convertThrottleDevice(spec.BlkioConfig.DeviceReadBPS),
			BlkioDeviceReadIOps:  convertThrottleDevice(spec.BlkioConfig.DeviceReadIOPS),
			BlkioDeviceWriteBps:  convertThrottleDevice(spec.BlkioConfig.DeviceWriteBPS),
			BlkioDeviceWriteIOps: convertThrottleDevice(spec.BlkioConfig.DeviceWriteIOPS),
			CpusetCpus:           spec.CPUSet,
			CgroupParent:         spec.CgroupParent,
			DeviceCgroupRules:    spec.DeviceCgroupRules,
			Devices:              convertDeviceMappings(spec.Devices),
			OomKillDisable:       spec.OomKillDisable,
			PidsLimit:            spec.PidsLimit,
			Ulimits:              convertUlimits(spec.Ulimits),
		},

		OomScoreAdj: spec.OomScoreAdj,

		//// MaskedPaths is the list of paths to be masked inside the container (this overrides the default set of paths)
		//MaskedPaths []string

		//// ReadonlyPaths is the list of paths to be set as read-only inside the container (this overrides the default set of paths)
		//ReadonlyPaths []string

		//// Run a custom init inside the container, if null, use the daemon's configured settings
		Init: spec.Init,
	}

	var err error
	if hostCfg.IpcMode, err = c.parseIPCMode(spec.IPC); err != nil {
		return err
	}

	if hostCfg.Isolation, err = parseIsolation(spec.Isolation); err != nil {
		return err
	}

	if hostCfg.NetworkMode, err = c.parseNetworkMode(spec.NetworkMode); err != nil {
		return err
	}

	if hostCfg.PidMode, err = c.parsePIDMode(spec.PidMode); err != nil {
		return err
	}

	if hostCfg.VolumesFrom, err = c.parseVolumesFrom(spec.VolumesFrom); err != nil {
		return err
	}

	if len(spec.Tmpfs) > 0 {
		hostCfg.Tmpfs = make(map[string]string, len(spec.Tmpfs))
		// This matches the docker-compose behaviour for specifying tmpfs mounts with the service-level `tmpfs` option.
		for _, path := range spec.Tmpfs {
			hostCfg.Tmpfs[path] = ""
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
		target, err := nat.NewPort(mapping.Protocol, mapping.Target)
		if err != nil {
			return fmt.Errorf("could not parse port: %w", err)
		}

		targetLow, targetHigh, err := target.Range()
		if err != nil {
			return fmt.Errorf("could not parse port range: %w", err)
		}

		for targetPort := targetLow; targetPort <= targetHigh; targetPort += 1 {
			publishedPort, err := nat.NewPort(mapping.Protocol, mapping.Published)
			if err != nil {
				return fmt.Errorf("could not parse port range: %w", err)
			}

			publishedLow, publishedHigh, err := publishedPort.Range()
			if err != nil {
				return fmt.Errorf("could not parse port range: %w", err)
			}

			publishedDiff, targetDiff := publishedHigh-publishedLow, targetHigh-targetLow
			if publishedDiff != 1 && publishedDiff != targetDiff {
				return fmt.Errorf("unexpected number of ports")
			}

			hostPort := publishedPort.Port()
			if publishedHigh != publishedLow {
				hostPort = strconv.Itoa(publishedLow + targetPort - targetLow)
			}

			bindings := hostCfg.PortBindings[target]
			bindings = append(bindings, nat.PortBinding{
				HostIP:   mapping.HostIP,
				HostPort: hostPort,
			})

			// TODO: Handle mapping.Mode
			hostCfg.PortBindings[nat.Port(strconv.Itoa(targetPort))] = bindings
		}
	}
	networkCfg := &network.NetworkingConfig{
		EndpointsConfig: make(map[string]*network.EndpointSettings), // Endpoint configs for each connecting network
	}
	// Docker only allows a single network to be specified when creating a container. The other networks must be
	// connected after the container is started. See https://github.com/moby/moby/issues/29265#issuecomment-265909198.
	var remainingNetworks []compose.ServiceNetwork
	if len(spec.Networks) > 0 {
		firstNetwork := spec.Networks[0]
		remainingNetworks = spec.Networks[1:]
		networkCfg.EndpointsConfig[firstNetwork.Network] = c.endpointSettings(firstNetwork, spec)
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
	createdBody, err := c.Docker.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, platform, spec.ContainerName)
	if err != nil {
		return err
	}
	c.State.ContainerID = createdBody.ID
	var netConnects errgroup.Group
	for _, network := range remainingNetworks {
		network := network
		netConnects.Go(func() error {
			return c.Docker.NetworkConnect(ctx, network.Network, createdBody.ID, c.endpointSettings(network, spec))
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
		if err := yamlutil.UnmarshalString(input.Spec, &spec); err != nil {
			return nil, fmt.Errorf("unmarshalling spec: %w", err)
		}
		c.State.Image.Spec = yamlutil.MustMarshalString(image.Spec{
			Platform: spec.Platform,
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
		Aliases: append([]string{c.ComponentID, c.ComponentName}, sn.Aliases...),
		Links:   append(append([]string{}, spec.Links...), spec.ExternalLinks...),
	}

	if sn.IPV4Address != "" || sn.IPV6Address != "" || len(sn.LinkLocalIPs) > 0 {
		es.IPAMConfig = &network.EndpointIPAMConfig{
			IPv4Address:  sn.IPV4Address,
			IPv6Address:  sn.IPV6Address,
			LinkLocalIPs: sn.LinkLocalIPs,
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
			Path: throttleDevice.Path,
			Rate: uint64(throttleDevice.Rate),
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
			Hard: ulimit.Hard,
			Soft: ulimit.Soft,
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
