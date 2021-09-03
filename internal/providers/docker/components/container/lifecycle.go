package container

import (
	"context"
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/blkiodev"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
)

var _ core.Lifecycle = (*Container)(nil)

func (c *Container) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {

	if err := c.ensureImage(ctx); err != nil {
		return nil, fmt.Errorf("ensuring image: %w", err)
	}

	if err := c.removeExistingContainerByName(ctx); err != nil {
		return nil, fmt.Errorf("removing existing container %q: %w", c.Spec.ContainerName, err)
	}

	if err := c.create(ctx); err != nil {
		return nil, fmt.Errorf("creating container: %w", err)
	}

	if err := c.start(ctx); err != nil {
		c.Logger.Infof("starting container %q: %v", c.State.ContainerID, err)
	}

	return &core.InitializeOutput{}, nil
}

func (c *Container) create(ctx context.Context) error {
	var healthCfg *container.HealthConfig
	if c.Spec.Healthcheck != nil {
		healthCfg = &container.HealthConfig{
			Test:        strslice.StrSlice(c.Spec.Healthcheck.Test.Parts),
			Interval:    time.Duration(c.Spec.Healthcheck.Interval),
			Timeout:     time.Duration(c.Spec.Healthcheck.Timeout),
			Retries:     c.Spec.Healthcheck.Retries,
			StartPeriod: time.Duration(c.Spec.Healthcheck.StartPeriod),
		}
	}

	labels := c.Spec.Labels.WithoutNils()
	for k, v := range c.GetExoLabels() {
		labels[k] = v
	}

	envMap := map[string]string{}
	for k, v := range c.Spec.Environment {
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
		Hostname:     c.Spec.Hostname,
		Domainname:   c.Spec.Domainname,
		User:         c.Spec.User,
		ExposedPorts: make(nat.PortSet),
		Tty:          c.Spec.TTY,
		OpenStdin:    c.Spec.StdinOpen,
		// StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
		Env:         envSlice,
		Healthcheck: healthCfg,
		// ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (meaning treat as a command line) (Windows specific).

		Image: c.State.Image.ID,
		// Volumes         map[string]struct{} // List of volumes (mounts) used for the container
		WorkingDir: c.Spec.WorkingDir,
		Entrypoint: strslice.StrSlice(c.Spec.Entrypoint.Parts),
		// NetworkDisabled bool                `json:",omitempty"` // Is network disabled
		MacAddress: c.Spec.MacAddress,
		// OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
		Labels:     labels,
		StopSignal: c.Spec.StopSignal,
		// Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}

	if c.Spec.Command.IsShellForm {
		containerCfg.Cmd = append(c.State.Image.Shell, c.Spec.Command.Parts[0])
	} else {
		containerCfg.Cmd = c.Spec.Command.Parts
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

	if c.Spec.StopGracePeriod != nil {
		timeout := int(time.Duration(*c.Spec.StopGracePeriod).Round(time.Second).Seconds())
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

	for _, exposed := range c.Spec.Expose {
		if err := exposePort(exposed.Target, exposed.Protocol); err != nil {
			return fmt.Errorf("exposing port %q: %w", exposed.Target, err)
		}
	}
	for _, mapping := range c.Spec.Ports {
		if err := exposePort(mapping.Target, mapping.Protocol); err != nil {
			return fmt.Errorf("exposing mapped port %q: %w", mapping.Target, err)
		}
	}

	logCfg := container.LogConfig{}
	if c.Spec.Logging.Driver == "" && (c.Spec.Logging.Options == nil || len(c.Spec.Logging.Options) == 0) {
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
		logCfg.Type = c.Spec.Logging.Driver
		logCfg.Config = c.Spec.Logging.Options
	}

	blkioWeightDevice := make([]*blkiodev.WeightDevice, len(c.Spec.BlkioConfig.WeightDevice))
	for i, weightDevice := range c.Spec.BlkioConfig.WeightDevice {
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
		//RestartPolicy   RestartPolicy // Restart policy to be used for the container
		// TODO: Potentially inherit from deploy's restart_policy.
		RestartPolicy: container.RestartPolicy{
			Name: c.Spec.Restart,
		},
		//AutoRemove      bool          // Automatically remove container when it exits
		//VolumeDriver    string        // Name of the volume driver used to mount volumes
		//VolumesFrom     []string      // List of volumes to take from other container

		//// Applicable to UNIX platforms
		CapAdd:  c.Spec.CapAdd,
		CapDrop: c.Spec.CapDrop,
		//CgroupnsMode    CgroupnsMode      // Cgroup namespace mode to use for the container
		DNS:        c.Spec.DNS,
		DNSOptions: c.Spec.DNSOptions,
		DNSSearch:  c.Spec.DNSSearch,
		ExtraHosts: c.Spec.ExtraHosts,
		GroupAdd:   c.Spec.GroupAdd,
		//IpcMode         IpcMode           // IPC namespace to use for the container
		//Cgroup          CgroupSpec        // Cgroup to use for the container

		// See NOTE: [RESOLVING SERVICE CONTAINERS].
		Links: append(append([]string{}, c.Spec.Links...), c.Spec.ExternalLinks...),
		//OomScoreAdj     int               // Container preference for OOM-killing
		//PidMode         PidMode           // PID namespace to use for the container
		Privileged: c.Spec.Privileged,
		//PublishAllPorts bool              // Should docker publish all exposed port for the container
		//ReadonlyRootfs  bool              // Is the container root filesystem in read-only
		SecurityOpt: c.Spec.SecurityOpt,
		//StorageOpt      map[string]string `json:",omitempty"` // Storage driver options per container.
		//Tmpfs           map[string]string `json:",omitempty"` // List of tmpfs (mounts) used for the container
		//UTSMode         UTSMode           // UTS namespace to use for the container
		//UsernsMode      UsernsMode        // The user namespace to use for the container
		ShmSize: int64(c.Spec.ShmSize),
		//Sysctls         map[string]string `json:",omitempty"` // List of Namespaced sysctls used for the container
		Runtime: c.Spec.Runtime,

		//// Applicable to Windows
		//ConsoleSize [2]uint   // Initial console size (height,width)

		//// Contains container's resources (cgroups, ulimits)
		Resources: container.Resources{
			CPUCount:             c.Spec.CPUCount,
			CPUPercent:           c.Spec.CPUPercent,
			CPUShares:            c.Spec.CPUShares,
			CPUPeriod:            c.Spec.CPUPeriod,
			CPUQuota:             c.Spec.CPUQuota,
			Memory:               int64(c.Spec.MemoryLimit),
			MemoryReservation:    int64(c.Spec.MemoryReservation),
			MemorySwappiness:     c.Spec.MemorySwappiness,
			CPURealtimePeriod:    time.Duration(c.Spec.CPURealtimePeriod).Microseconds(),
			CPURealtimeRuntime:   time.Duration(c.Spec.CPURealtimeRuntime).Microseconds(),
			BlkioWeight:          uint16(c.Spec.BlkioConfig.Weight),
			BlkioWeightDevice:    blkioWeightDevice,
			BlkioDeviceReadBps:   convertThrottleDevice(c.Spec.BlkioConfig.DeviceReadBPS),
			BlkioDeviceReadIOps:  convertThrottleDevice(c.Spec.BlkioConfig.DeviceReadIOPS),
			BlkioDeviceWriteBps:  convertThrottleDevice(c.Spec.BlkioConfig.DeviceWriteBPS),
			BlkioDeviceWriteIOps: convertThrottleDevice(c.Spec.BlkioConfig.DeviceWriteIOPS),
			CpusetCpus:           c.Spec.CPUSet,
			CgroupParent:         c.Spec.CgroupParent,
			DeviceCgroupRules:    c.Spec.DeviceCgroupRules,
			Devices:              convertDeviceMappings(c.Spec.Devices),
		},

		// Mounts specs used by the container
		//Mounts []mount.Mount `json:",omitempty"`

		//// MaskedPaths is the list of paths to be masked inside the container (this overrides the default set of paths)
		//MaskedPaths []string

		//// ReadonlyPaths is the list of paths to be set as read-only inside the container (this overrides the default set of paths)
		//ReadonlyPaths []string

		//// Run a custom init inside the container, if null, use the daemon's configured settings
		Init: c.Spec.Init,
	}

	var err error
	if hostCfg.IpcMode, err = c.parseIPCMode(c.Spec.IPC); err != nil {
		return err
	}

	if hostCfg.Isolation, err = parseIsolation(c.Spec.Isolation); err != nil {
		return err
	}

	// TODO: make the user home directory a parameter of the container.
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not get user %w", err)
	}
	userHomeDir := user.HomeDir

	hostCfg.Mounts = make([]mount.Mount, len(c.Spec.Volumes))
	for i, v := range c.Spec.Volumes {
		mnt, err := makeMountFromVolumeMount(c.WorkspaceRoot, userHomeDir, v)
		if err != nil {
			return fmt.Errorf("invalid mount at index %d: %w", i, err)
		}
		hostCfg.Mounts[i] = mnt
	}

	for _, mapping := range c.Spec.Ports {
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
	netEndpointSettings := &network.EndpointSettings{
		// TODO: Add other specified aliases.
		Aliases: []string{c.ComponentID, c.ComponentName},
	}
	// Docker only allows a single network to be specified when creating a container. The other networks must be
	// connected after the container is started. See https://github.com/moby/moby/issues/29265#issuecomment-265909198.
	var remainingNetworks []string
	if len(c.Spec.Networks) > 0 {
		firstNetworkName := c.Spec.Networks[0]
		remainingNetworks = c.Spec.Networks[1:]
		networkCfg.EndpointsConfig[firstNetworkName] = netEndpointSettings
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
	createdBody, err := c.Docker.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, platform, c.Spec.ContainerName)
	if err != nil {
		return err
	}
	c.State.ContainerID = createdBody.ID
	var netConnects errgroup.Group
	for _, networkName := range remainingNetworks {
		networkName := networkName
		netConnects.Go(func() error {
			return c.Docker.NetworkConnect(ctx, networkName, createdBody.ID, netEndpointSettings)
		})
	}

	return netConnects.Wait()
}

func (c *Container) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	if c.State.ContainerID == "" {
		c.State.Running = false
		return &core.RefreshOutput{}, nil
	}

	inspection, err := c.Docker.ContainerInspect(ctx, c.State.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("inspecting container: %w", err)
	}

	c.State.Running = inspection.State.Running
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
		c.Logger.Infof("disposing container not found: %q", c.State.ContainerID)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	c.State.ContainerID = ""
	return &core.DisposeOutput{}, nil
}

func (c *Container) removeExistingContainerByName(ctx context.Context) error {
	// If a container with this name already exists, remove it.
	containers, err := c.Docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: c.Spec.ContainerName,
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
