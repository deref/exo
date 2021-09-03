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
	"github.com/goccy/go-yaml"
)

func Parse(r io.Reader) (*Compose, error) {
	dec := yaml.NewDecoder(r,
		yaml.DisallowDuplicateKey(),
	)
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
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
	Volumes  map[string]Volume  `yaml:"volumes"`
	Configs  map[string]Config  `yaml:"configs"`
	Secrets  map[string]Secret  `yaml:"secrets"`

	Raw map[string]interface{} `yaml:",inline"`
	// TODO: extensions with "x-" prefix.
}

// This is a temporary placeholder for fields that we presently don't support,
// but are safe to ignore.
// TODO: Eliminate all usages of this with actual parsing logic.
type IgnoredField struct{}

func (ignored *IgnoredField) UnmarshalYAML(b []byte) error {
	return nil
}

type MemoryField int64

func (memory *MemoryField) UnmarshalYAML(b []byte) error {
	memString := string(b)
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

	return fmt.Errorf("could not unmarshal memory value %s: %w", b, err)
}

type Service struct {
	Deploy IgnoredField `yaml:"deploy"`

	// Note that these two are only applicable to Windows.
	CPUCount   int64 `yaml:"cpu_count"`
	CPUPercent int64 `yaml:"cpu_percent"`

	CPUShares int64 `yaml:"cpu_shares"`
	CPUPeriod int64 `yaml:"cpu_period"`
	CPUQuota  int64 `yaml:"cpu_quota"`

	CPURealtimeRuntime Duration `yaml:"cpu_rt_runtime"`
	CPURealtimePeriod  Duration `yaml:"cpu_rt_period"`

	CPUSet string `yaml:"cpuset"`

	// TODO: issue deprecation warning if `cpus` is set (replaced by `deploy.reservations.cpus`)

	BlkioConfig BlkioConfig `yaml:"blkio_config"`

	Build Build `yaml:"build"`

	CapAdd  []string `yaml:"cap_add"`
	CapDrop []string `yaml:"cap_drop"`

	CgroupParent string `yaml:"cgroup_parent"`

	Command       Command  `yaml:"command"`
	Configs       []string `yaml:"configs"` // TODO: support long syntax.
	ContainerName string   `yaml:"container_name"`
	// TODO: credential_spec

	DependsOn ServiceDependencies `yaml:"depends_on"`

	DeviceCgroupRules []string `yaml:"device_cgroup_rules"`

	Devices     []DeviceMapping     `yaml:"devices"`
	DNS         StringOrStringSlice `yaml:"dns"`
	DNSOptions  []string            `yaml:"dns_opt"`
	DNSSearch   StringOrStringSlice `yaml:"dns_search"`
	Domainname  string              `yaml:"domainname"`
	Entrypoint  Command             `yaml:"entrypoint"`
	EnvFile     StringOrStringSlice `yaml:"env_file"` // TODO: Add to the environment for a docker container component.
	Environment Dictionary          `yaml:"environment"`
	Expose      PortMappings        `yaml:"expose"` // TODO: Validate target-only.
	// TODO: extends

	// List of links of the form `SERVICE` or `SERVICE:ALIAS`
	ExternalLinks []string `yaml:"external_links"`

	// List of host/IP pairs to add to /etc/hosts of the form `HOST:IP`
	ExtraHosts []string `yaml:"extra_hosts"`

	GroupAdd    []string     `yaml:"group_add"`
	Healthcheck *Healthcheck `yaml:"healthcheck"`
	Hostname    string       `yaml:"hostname"`
	Image       string       `yaml:"image"`
	Init        *bool        `yaml:"init"`
	// TODO: init
	// TODO: ipc
	// TODO: isolation
	Labels Dictionary `yaml:"labels"`
	// TODO: links
	Logging Logging `yaml:"logging"`
	// TODO: network_mode
	Networks   []string `yaml:"networks"` // TODO: support long syntax.
	MacAddress string   `yaml:"mac_address"`

	MemorySwappiness *int64 `yaml:"mem_swappiness"`

	// MemoryLimit and MemoryReservation can be specified either as strings or integers.
	MemoryLimit       MemoryField `yaml:"mem_limit"`
	MemoryReservation MemoryField `yaml:"mem_reservation"`

	// TODO: memswap_limit
	// TODO: oom_kill_disable
	// TODO: oom_score_adj
	// TODO: pid
	// TODO: pids_limit
	// TODO: platform

	Ports      PortMappings `yaml:"ports"`
	Privileged bool         `yaml:"privileged"`
	Profiles   IgnoredField `yaml:"profiles"`
	// TODO: pull_policy
	// TODO: read_only
	Restart string `yaml:"restart"`
	Runtime string `yaml:"runtime"`
	// TODO: scale
	Secrets         []string  `yaml:"secrets"` // TODO: support long syntax.
	SecurityOpt     []string  `yaml:"security_opt"`
	ShmSize         Bytes     `yaml:"shm_size"`
	StdinOpen       bool      `yaml:"stdin_open"`
	StopGracePeriod *Duration `yaml:"stop_grace_period"`
	StopSignal      string    `yaml:"stop_signal"`
	// TODO: storage_opt
	// TODO: sysctls
	// TODO: tmpfs
	TTY bool `yaml:"tty"`
	// TODO: ulimits
	User string `yaml:"user"`
	// TODO: userns_mode
	Volumes []VolumeMount `yaml:"volumes"`
	// TODO: volumes_from

	WorkingDir string `yaml:"working_dir"`
}

type Healthcheck struct {
	Test        Command  `yaml:"test"`
	Interval    Duration `yaml:"interval"`
	Timeout     Duration `yaml:"timeout"`
	Retries     int      `yaml:"retries"`
	StartPeriod Duration `yaml:"start_period"`
}

type Logging struct {
	Driver  string            `yaml:"driver"`
	Options map[string]string `yaml:"options"`
}

type Network struct {
	// Name is the actual name of the docker network. The docker-compose network name, which can
	// be referenced by individual services, is the component name.
	Name       string            `yaml:"name"`
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	Attachable bool              `yaml:"attachable"`
	EnableIPv6 bool              `yaml:"enable_ipv6"`
	Internal   bool              `yaml:"internal"`
	Labels     Dictionary        `yaml:"labels"`
	External   bool              `yaml:"external"`
}

type Volume struct {
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	// TODO: external
	Labels Dictionary `yaml:"labels"`
	Name   string     `yaml:"name"`
}

type Config struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}

type Secret struct {
	File     string `yaml:"file"`
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
}

type BlkioConfig struct {
	DeviceReadBPS   []ThrottleDevice `yaml:"device_read_bps"`
	DeviceWriteBPS  []ThrottleDevice `yaml:"device_write_bps"`
	DeviceReadIOPS  []ThrottleDevice `yaml:"device_read_iops"`
	DeviceWriteIOPS []ThrottleDevice `yaml:"device_write_iops"`
	Weight          uint16           `yaml:"weight"`
	WeightDevice    []WeightDevice   `yaml:"weight_device"`
}

type ThrottleDevice struct {
	Path string `yaml:"path"`
	Rate Bytes  `yaml:"rate"`
}

type WeightDevice struct {
	Path   string `yaml:"path"`
	Weight uint16 `yaml:"weight"`
}
