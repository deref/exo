// References:
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// https://docs.docker.com/compose/compose-file/compose-file-v3/
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
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
	Volumes  map[string]Volume  `yaml:"volumes"`
	Configs  map[string]Config  `yaml:"configs"`
	Secrets  map[string]Secret  `yaml:"secrets"`
	// TODO: extensions with "x-" prefix.
}

type Service struct {
	// XXX fill me.
	// TODO: deploy
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
	// TODO: command
	Configs       []string `yaml:"configs"` // TODO: support long syntax.
	ContainerName string   `yaml:"container_name"`
	// TODO: credential_spec
	// TODO: depends_on
	// TODO: device_cgroup_rules
	// TODO: devices
	// TODO: dns
	// TODO: dns_opt
	// TODO: dns_search
	// TODO: domainname
	// TODO: entrypoint
	// TODO: env_file
	Environment Dictionary `yaml:"environment"`
	// TODO: expose
	// TODO: extends
	// TODO: external_links
	// TODO: extra_hosts
	// TODO: group_add
	// TODO: healthcheck
	// TODO: hostname
	Image string `yaml:"image"`
	// TODO: init
	// TODO: ipc
	// TODO: isolation
	// TODO: labels
	// TODO: links
	Logging Logging `yaml:"logging"`
	// TODO: network_mode
	Networks []string `yaml:"networks"` // TODO: support long syntax.
	// TODO: mac_address
	// TODO: mem_limit
	// TODO: mem_reservation
	// TODO: mem_swappiness
	// TODO: memswap_limit
	// TODO: oom_kill_disable
	// TODO: oom_score_adj
	// TODO: pid
	// TODO: pids_limit
	// TODO: platform
	Ports PortMappings `yaml:"ports"`
	// TODO: privileged
	// TODO: profiles
	// TODO: pull_policy
	// TODO: read_only
	Restart string `yaml:"restart"`
	// TODO: runtime
	// TODO: scale
	Secrets []string `yaml:"secrets"` // TODO: support long syntax.
	// TODO: security_opt
	// TODO: shm_size
	// TODO: shm_open
	// TODO: stop_grace_period
	// TODO: stop_signal
	// TODO: storage_opt
	// TODO: sysctls
	// TODO: tmpfs
	// TODO: tty
	// TODO: ulimits
	// TODO: user
	// TODO: userns_mode
	Volumes []string `yaml:"volumes"` // TODO: support long syntax.
	// TODO: volumes_from
	// TODO: working_dir
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
