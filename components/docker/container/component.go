package container

import (
	docker "github.com/docker/docker/client"
)

type Container struct {
	ComponentID string
	Spec
	State

	Docker *docker.Client
}

// See note: [COMPOSE_YAML].
type Spec struct {
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
	// TODO: build
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
	Environment []string `yaml:"environment"` // TODO: Support map syntax.
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
	// TODO: logging
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
	Ports []string `yaml:"ports"` // TODO: support long syntax.
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

type State struct {
	ContainerID string `json:"containerId"`
}
