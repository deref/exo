package compose

type Service struct {
	Key string `yaml:"-"`

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
	DependsOn         ServiceDependencies     `yaml:"depends_on,omitempty"`
	DeviceCgroupRules []string                `yaml:"device_cgroup_rules,omitempty"`
	Devices           []DeviceMapping         `yaml:"devices,omitempty"`
	DNS               Strings                 `yaml:"dns,omitempty"`
	DNSOptions        []string                `yaml:"dns_opt,omitempty"`
	DNSSearch         Strings                 `yaml:"dns_search,omitempty"`
	Domainname        string                  `yaml:"domainname,omitempty"`
	Entrypoint        Command                 `yaml:"entrypoint,omitempty"`
	EnvFile           Strings                 `yaml:"env_file,omitempty"`
	Environment       Dictionary              `yaml:"environment,omitempty"`
	Expose            []PortRangeWithProtocol `yaml:"expose,omitempty"`
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
	MemoryLimit       Bytes `yaml:"mem_limit,omitempty"`
	MemoryReservation Bytes `yaml:"mem_reservation,omitempty"`

	MemswapLimit   Bytes        `yaml:"memswap_limit,omitempty"`
	OomKillDisable *bool        `yaml:"oom_kill_disable,omitempty"`
	OomScoreAdj    int          `yaml:"oom_score_adj,omitempty"`
	PidMode        string       `yaml:"pid,omitempty"`
	PidsLimit      *int64       `yaml:"pids_limit,omitempty"`
	Platform       string       `yaml:"platform,omitempty"`
	Ports          PortMappings `yaml:"ports,omitempty"`
	Privileged     bool         `yaml:"privileged,omitempty"`
	// TODO: Support profiles. See https://docs.docker.com/compose/profiles/.
	Profiles        Ignored       `yaml:"profiles,omitempty"`
	PullPolicy      string        `yaml:"pull_policy,omitempty"`
	ReadOnly        bool          `yaml:"read_only,omitempty"`
	Restart         string        `yaml:"restart,omitempty"`
	Runtime         string        `yaml:"runtime,omitempty"`
	SecurityOpt     []string      `yaml:"security_opt,omitempty"`
	ShmSize         Bytes         `yaml:"shm_size,omitempty"`
	StdinOpen       bool          `yaml:"stdin_open,omitempty"`
	StopGracePeriod *Duration     `yaml:"stop_grace_period,omitempty"`
	StopSignal      string        `yaml:"stop_signal,omitempty"`
	StorageOpt      Dictionary    `yaml:"storage_opt,omitempty"`
	Sysctls         Dictionary    `yaml:"sysctls,omitempty"`
	Tmpfs           Strings       `yaml:"tmpfs,omitempty"`
	TTY             bool          `yaml:"tty,omitempty"`
	Ulimits         Ulimits       `yaml:"ulimits,omitempty"`
	User            string        `yaml:"user,omitempty"`
	UsernsMode      string        `yaml:"userns_mode,omitempty"`
	Volumes         []VolumeMount `yaml:"volumes,omitempty"`
	VolumesFrom     []string      `yaml:"volumes_from,omitempty"`
	WorkingDir      string        `yaml:"working_dir,omitempty"`

	// NOTE [DOCKER SWARM FEATURES]:
	// Docker-Compose manages local, single-container deployments as well as Docker Swarm
	// deployments. Since Swarm is not as widely used as Kubernetes, support for the Swarm
	// features that Docker-Compose includes is not a top priority. The settings listed
	// below are the ones that are applicable to a Swarm deployment.
	Deploy  Ignored `yaml:"deploy,omitempty"`
	Scale   Ignored `yaml:"scale,omitempty"`
	Secrets Ignored `yaml:"secrets,omitempty"`
}
