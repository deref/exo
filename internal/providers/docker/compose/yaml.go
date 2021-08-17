// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
// https://github.com/docker/compose/blob/4a51af09d6cdb9407a6717334333900327bc9302/compose/config/compose_spec.json
//
// Fields enumerated as of July 17, 2021 with from the following spec file:
// <https://github.com/compose-spec/compose-spec/blob/5141aafafa6ea03fcf52eb2b44218408825ab480/spec.md>.

package compose

import (
	"io"

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
	Version  String             `yaml:"version"`
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

type Service struct {
	Deploy IgnoredField `yaml:"deploy"`

	// TODO: blkio_config
	// TODO: cpu_count
	// TODO: cpu_percent
	// TODO: cpu_shares
	// TODO: cpu_period
	// TODO: cpu_quota
	// TODO: cpu_rt_runtime
	// TODO: cpu_rt_period
	// TODO: cpus
	// TODO: cpuset
	Build Build `yaml:"build"`
	// TODO: cap_add
	// TODO: cap_drop
	// TODO: cgroup_parent

	Command       Command  `yaml:"command"`
	Configs       []String `yaml:"configs"` // TODO: support long syntax.
	ContainerName String   `yaml:"container_name"`
	// TODO: credential_spec

	DependsOn IgnoredField `yaml:"depends_on"`

	// TODO: device_cgroup_rules
	// TODO: devices
	// TODO: dns
	// TODO: dns_opt
	// TODO: dns_search
	Domainname String  `yaml:"domainname"`
	Entrypoint Command `yaml:"entrypoint"`
	// TODO: env_file
	Environment Dictionary   `yaml:"environment"`
	Expose      PortMappings `yaml:"expose"` // TODO: Validate target-only.
	// TODO: extends
	// TODO: external_links
	// TODO: extra_hosts
	// TODO: group_add
	Healthcheck *Healthcheck `yaml:"healthcheck"`
	Hostname    String       `yaml:"hostname"`
	Image       String       `yaml:"image"`
	// TODO: init
	// TODO: ipc
	// TODO: isolation
	Labels Dictionary `yaml:"labels"`
	// TODO: links
	Logging Logging `yaml:"logging"`
	// TODO: network_mode
	Networks   []String `yaml:"networks"` // TODO: support long syntax.
	MacAddress String   `yaml:"mac_address"`
	// TODO: mem_limit
	// TODO: mem_reservation
	// TODO: mem_swappiness
	// TODO: memswap_limit
	// TODO: oom_kill_disable
	// TODO: oom_score_adj
	// TODO: pid
	// TODO: pids_limit
	// TODO: platform
	Ports      PortMappings `yaml:"ports"`
	Privileged Bool         `yaml:"privileged"`
	Profiles   IgnoredField `yaml:"profiles"`
	// TODO: pull_policy
	// TODO: read_only
	Restart String `yaml:"restart"`
	Runtime String `yaml:"runtime"`
	// TODO: scale
	Secrets []String `yaml:"secrets"` // TODO: support long syntax.
	// TODO: security_opt
	ShmSize         Bytes     `yaml:"shm_size"`
	StdinOpen       Bool      `yaml:"stdin_open"`
	StopGracePeriod *Duration `yaml:"stop_grace_period"`
	StopSignal      String    `yaml:"stop_signal"`
	// TODO: storage_opt
	// TODO: sysctls
	// TODO: tmpfs
	TTY Bool `yaml:"tty"`
	// TODO: ulimits
	User String `yaml:"user"`
	// TODO: userns_mode
	Volumes []String `yaml:"volumes"` // TODO: support long syntax.
	// TODO: volumes_from

	WorkingDir String `yaml:"working_dir"`
}

type Healthcheck struct {
	Test        Command  `yaml:"test"`
	Interval    Duration `yaml:"interval"`
	Timeout     Duration `yaml:"timeout"`
	Retries     Int      `yaml:"retries"`
	StartPeriod Duration `yaml:"start_period"`
}

type Logging struct {
	Driver  String            `yaml:"driver"`
	Options map[string]String `yaml:"options"`
}

type Network struct {
	Driver     String            `yaml:"driver"`
	DriverOpts map[string]String `yaml:"driver_opts"`
	Attachable Bool              `yaml:"attachable"`
	EnableIPv6 Bool              `yaml:"enable_ipv6"`
	Internal   Bool              `yaml:"internal"`
	Labels     Dictionary        `yaml:"labels"`
	External   Bool              `yaml:"external"`
	// TODO: name
}

type Volume struct {
	Driver     String            `yaml:"driver"`
	DriverOpts map[string]String `yaml:"driver_opts"`
	// TODO: external
	Labels Dictionary `yaml:"labels"`
	Name   String     `yaml:"name"`
}

type Config struct {
	File     String `yaml:"file"`
	External Bool   `yaml:"external"`
	Name     String `yaml:"name"`
}

type Secret struct {
	File     String `yaml:"file"`
	External Bool   `yaml:"external"`
	Name     String `yaml:"name"`
}
