// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
// https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/compose_spec.json
//
// Fields enumerated as of July 17, 2021 with from the following spec file:
// <https://github.com/compose-spec/compose-spec/blob/5141aafafa6ea03fcf52eb2b44218408825ab480/spec.md>.

package compose

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"gopkg.in/yaml.v3"
)

func Parse(r io.Reader) (*Compose, error) {
	dec := yaml.NewDecoder(r)
	var comp Compose
	if err := dec.Decode(&comp); err != nil {
		return nil, err
	}

	// Validate.
	for key := range comp.Raw {
		switch key {
		case "version", "services", "networks", "volumes", "configs", "secrets":
			// Ok.
		default:
			if !strings.HasPrefix(key, "x-") {
				return nil, fmt.Errorf("unsupported top-level key in compose file: %q", key)
			}
		}
	}

	return &comp, nil
}

type Compose struct {
	Version  string             `yaml:"version,omitempty"`
	Services map[string]Service `yaml:"services,omitempty"`
	Networks map[string]Network `yaml:"networks,omitempty"`
	Volumes  map[string]Volume  `yaml:"volumes,omitempty"`
	Configs  map[string]Config  `yaml:"configs,omitempty"`
	Secrets  map[string]Secret  `yaml:"secrets,omitempty"`

	Raw map[string]interface{} `yaml:",inline"`
	// TODO: extensions with "x-" prefix.
}

// This is a temporary placeholder for fields that we presently don't support,
// but are safe to ignore.
// TODO: Eliminate all usages of this with actual parsing logic.
type IgnoredField struct{}

func (ignored *IgnoredField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return nil
}

type MemoryField int64

func (memory *MemoryField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var memString string
	if err := unmarshal(&memString); err != nil {
		return err
	}
	memBytes, err := strconv.ParseInt(memString, 10, 64)
	if err == nil {
		*memory = MemoryField(memBytes)
		return nil
	}

	uMemBytes, err := bytefmt.ToBytes(memString)
	if err == nil {
		*memory = MemoryField(uMemBytes)
		return nil
	}

	return fmt.Errorf("unmarshaling memory value: %w", err)
}

type Service struct {
	// Note that these two are only applicable to Windows.
	CPUCount   int64 `yaml:"cpu_count,omitempty"`
	CPUPercent int64 `yaml:"cpu_percent,omitempty"`

	CPUShares          int64       `yaml:"cpu_shares,omitempty"`
	CPUPeriod          int64       `yaml:"cpu_period,omitempty"`
	CPUQuota           int64       `yaml:"cpu_quota,omitempty"`
	CPURealtimeRuntime Duration    `yaml:"cpu_rt_runtime,omitempty"`
	CPURealtimePeriod  Duration    `yaml:"cpu_rt_period,omitempty"`
	CPUSet             string      `yaml:"cpuset,omitempty"`
	BlkioConfig        BlkioConfig `yaml:"blkio_config,omitempty"`
	Build              Build       `yaml:"build,omitempty"`
	CapAdd             []string    `yaml:"cap_add,omitempty"`
	CapDrop            []string    `yaml:"cap_drop,omitempty"`
	CgroupParent       string      `yaml:"cgroup_parent,omitempty"`
	Command            Command     `yaml:"command,omitempty"`
	Configs            []string    `yaml:"configs,omitempty"` // TODO: support long syntax.
	ContainerName      string      `yaml:"container_name,omitempty"`
	// TODO: credential_spec
	DependsOn         ServiceDependencies `yaml:"depends_on,omitempty"`
	DeviceCgroupRules []string            `yaml:"device_cgroup_rules,omitempty"`
	Devices           []DeviceMapping     `yaml:"devices,omitempty"`
	DNS               StringOrStringSlice `yaml:"dns,omitempty"`
	DNSOptions        []string            `yaml:"dns_opt,omitempty"`
	DNSSearch         StringOrStringSlice `yaml:"dns_search,omitempty"`
	Domainname        string              `yaml:"domainname,omitempty"`
	Entrypoint        Command             `yaml:"entrypoint,omitempty"`
	EnvFile           StringOrStringSlice `yaml:"env_file,omitempty"` // TODO: Add to the environment for a docker container component.
	Environment       Dictionary          `yaml:"environment,omitempty"`
	Expose            PortMappings        `yaml:"expose,omitempty"` // TODO: Validate target-only.
	// TODO: extends
	// List of links of the form `SERVICE` or `SERVICE:ALIAS`
	ExternalLinks []string `yaml:"external_links,omitempty"`
	// List of host/IP pairs to add to /etc/hosts of the form `HOST:IP`
	ExtraHosts       []string        `yaml:"extra_hosts,omitempty"`
	GroupAdd         []string        `yaml:"group_add,omitempty"`
	Healthcheck      *Healthcheck    `yaml:"healthcheck,omitempty"`
	Hostname         string          `yaml:"hostname,omitempty"`
	Image            string          `yaml:"image,omitempty"`
	Init             *bool           `yaml:"init,omitempty"`
	IPC              string          `yaml:"ipc,omitempty"`
	Isolation        string          `yaml:"isolation,omitempty"`
	Labels           Dictionary      `yaml:"labels,omitempty"`
	Links            []string        `yaml:"links,omitempty"`
	Logging          Logging         `yaml:"logging,omitempty"`
	NetworkMode      string          `yaml:"network_mode,omitempty"`
	Networks         ServiceNetworks `yaml:"networks,omitempty"`
	MacAddress       string          `yaml:"mac_address,omitempty"`
	MemorySwappiness *int64          `yaml:"mem_swappiness,omitempty"`
	// MemoryLimit and MemoryReservation can be specified either as strings or integers.
	// TODO: Deprecate these fields once we support `deploy.limits.memory` and `deploy.reservations.memory`.
	MemoryLimit       MemoryField `yaml:"mem_limit,omitempty"`
	MemoryReservation MemoryField `yaml:"mem_reservation,omitempty"`

	MemswapLimit   Bytes        `yaml:"memswap_limit,omitempty"`
	OomKillDisable *bool        `yaml:"oom_kill_disable,omitempty"`
	OomScoreAdj    int          `yaml:"oom_score_adj,omitempty"`
	PidMode        string       `yaml:"pid,omitempty"`
	PidsLimit      *int64       `yaml:"pids_limit,omitempty"`
	Platform       string       `yaml:"platform,omitempty"`
	Ports          PortMappings `yaml:"ports,omitempty"`
	Privileged     bool         `yaml:"privileged,omitempty"`
	// TODO: Support profiles. See https://docs.docker.com/compose/profiles/.
	Profiles        IgnoredField        `yaml:"profiles,omitempty"`
	PullPolicy      string              `yaml:"pull_policy,omitempty"`
	ReadOnly        bool                `yaml:"read_only,omitempty"`
	Restart         string              `yaml:"restart,omitempty"`
	Runtime         string              `yaml:"runtime,omitempty"`
	SecurityOpt     []string            `yaml:"security_opt,omitempty"`
	ShmSize         Bytes               `yaml:"shm_size,omitempty"`
	StdinOpen       bool                `yaml:"stdin_open,omitempty"`
	StopGracePeriod *Duration           `yaml:"stop_grace_period,omitempty"`
	StopSignal      string              `yaml:"stop_signal,omitempty"`
	StorageOpt      map[string]string   `yaml:"storage_opt,omitempty"`
	Sysctls         Dictionary          `yaml:"sysctls,omitempty"`
	Tmpfs           StringOrStringSlice `yaml:"tmpfs,omitempty"`
	TTY             bool                `yaml:"tty,omitempty"`
	Ulimits         Ulimits             `yaml:"ulimits,omitempty"`
	User            string              `yaml:"user,omitempty"`
	UsernsMode      string              `yaml:"userns_mode,omitempty"`
	Volumes         []VolumeMount       `yaml:"volumes,omitempty"`
	VolumesFrom     []string            `yaml:"volumes_from,omitempty"`
	WorkingDir      string              `yaml:"working_dir,omitempty"`

	// NOTE [DOCKER SWARM FEATURES]:
	// Docker-Compose manages local, single-container deployments as well as Docker Swarm
	// deployments. Since Swarm is not as widely used as Kubernetes, support for the Swarm
	// features that Docker-Compose includes is not a top priority. The settings listed
	// below are the ones that are applicable to a Swarm deployment.
	Deploy  IgnoredField `yaml:"deploy,omitempty"`
	Scale   IgnoredField `yaml:"scale,omitempty"`
	Secrets IgnoredField `yaml:"secrets,omitempty"`
}

type Healthcheck struct {
	Test        Command  `yaml:"test,omitempty"`
	Interval    Duration `yaml:"interval,omitempty"`
	Timeout     Duration `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod Duration `yaml:"start_period,omitempty"`
}

type Logging struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

type Network struct {
	// Name is the actual name of the docker network. The docker-compose network name, which can
	// be referenced by individual services, is the component name.
	Name       string            `yaml:"name,omitempty"`
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	Attachable bool              `yaml:"attachable,omitempty"`
	EnableIPv6 bool              `yaml:"enable_ipv6,omitempty"`
	Internal   bool              `yaml:"internal,omitempty"`
	Labels     Dictionary        `yaml:"labels,omitempty"`
	External   bool              `yaml:"external,omitempty"`
}

type Volume struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	// TODO: external
	Labels Dictionary `yaml:"labels,omitempty"`
	Name   string     `yaml:"name,omitempty"`
}

type Config struct {
	File     string `yaml:"file,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type Secret struct {
	File     string `yaml:"file,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type BlkioConfig struct {
	DeviceReadBPS   []ThrottleDevice `yaml:"device_read_bps,omitempty"`
	DeviceWriteBPS  []ThrottleDevice `yaml:"device_write_bps,omitempty"`
	DeviceReadIOPS  []ThrottleDevice `yaml:"device_read_iops,omitempty"`
	DeviceWriteIOPS []ThrottleDevice `yaml:"device_write_iops,omitempty"`
	Weight          uint16           `yaml:"weight,omitempty"`
	WeightDevice    []WeightDevice   `yaml:"weight_device,omitempty"`
}

type ThrottleDevice struct {
	Path string `yaml:"path,omitempty"`
	Rate Bytes  `yaml:"rate,omitempty"`
}

type WeightDevice struct {
	Path   string `yaml:"path,omitempty"`
	Weight uint16 `yaml:"weight,omitempty"`
}
