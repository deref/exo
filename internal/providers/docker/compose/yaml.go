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

	"code.cloudfoundry.org/bytefmt"
	"github.com/goccy/go-yaml"
)

func Parse(r io.Reader) (*Compose, error) {
	dec := yaml.NewDecoder(r,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(), // TODO: Handle this more gracefully.
	)
	var comp Compose
	if err := dec.Decode(&comp); err != nil {
		return nil, err
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
	// TODO: cpu_count
	// TODO: cpu_percent

	CPUShares int64 `yaml:"cpu_shares"`
	CPUPeriod int64 `yaml:"cpu_period"`
	CPUQuota  int64 `yaml:"cpu_quota"`

	// CPURealtimeRuntime and CPURealtimePeriod can both be specified as either
	// strings or integers.
	CPURealtimeRuntime int64 `yaml:"cpu_rt_runtime"`
	CPURealtimePeriod  int64 `yaml:"cpu_rt_period"`

	// TODO: cpus
	// TODO: cpuset
	// TODO: blkio_config

	Build Build `yaml:"build"`
	// TODO: cap_add
	// TODO: cap_drop
	// TODO: cgroup_parent

	Command       Command  `yaml:"command"`
	Configs       []string `yaml:"configs"` // TODO: support long syntax.
	ContainerName string   `yaml:"container_name"`
	// TODO: credential_spec

	DependsOn IgnoredField `yaml:"depends_on"`

	// TODO: device_cgroup_rules
	// TODO: devices
	// TODO: dns
	// TODO: dns_opt
	// TODO: dns_search
	Domainname string  `yaml:"domainname"`
	Entrypoint Command `yaml:"entrypoint"`
	// TODO: env_file
	Environment Dictionary   `yaml:"environment"`
	Expose      PortMappings `yaml:"expose"` // TODO: Validate target-only.
	// TODO: extends
	// TODO: external_links
	// TODO: extra_hosts
	// TODO: group_add
	Healthcheck *Healthcheck `yaml:"healthcheck"`
	Hostname    string       `yaml:"hostname"`
	Image       string       `yaml:"image"`
	// TODO: init
	// TODO: ipc
	// TODO: isolation
	Labels Dictionary `yaml:"labels"`
	// TODO: links
	Logging Logging `yaml:"logging"`
	// TODO: network_mode
	Networks   []string `yaml:"networks"` // TODO: support long syntax.
	MacAddress string   `yaml:"mac_address"`

	// The following memory values can be specified either as strings or integers.
	MemoryLimit       int64 `yaml:"mem_limit"`
	MemoryReservation int64 `yaml:"mem_reservation"`
	MemorySwappiness  int64 `yaml:"mem_swappiness"`

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
	Secrets []string `yaml:"secrets"` // TODO: support long syntax.
	// TODO: security_opt
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
	Volumes []string `yaml:"volumes"` // TODO: support long syntax.
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
	Driver     string            `yaml:"driver"`
	DriverOpts map[string]string `yaml:"driver_opts"`
	Attachable bool              `yaml:"attachable"`
	EnableIPv6 bool              `yaml:"enable_ipv6"`
	Internal   bool              `yaml:"internal"`
	Labels     Dictionary        `yaml:"labels"`
	External   bool              `yaml:"external"`
	// TODO: name
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
