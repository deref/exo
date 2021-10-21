package compose

type Service struct {
	Key string `yaml:"-"`

	// Note that these two are only applicable to Windows.
	CPUCount   Int `yaml:"cpu_count,omitempty"`
	CPUPercent Int `yaml:"cpu_percent,omitempty"`

	CPUShares          Int         `yaml:"cpu_shares,omitempty"`
	CPUPeriod          Int         `yaml:"cpu_period,omitempty"`
	CPUQuota           Int         `yaml:"cpu_quota,omitempty"`
	CPURealtimeRuntime Duration    `yaml:"cpu_rt_runtime,omitempty"`
	CPURealtimePeriod  Duration    `yaml:"cpu_rt_period,omitempty"`
	CPUSet             String      `yaml:"cpuset,omitempty"`
	BlkioConfig        BlkioConfig `yaml:"blkio_config,omitempty"`
	Build              Build       `yaml:"build,omitempty"`
	CapAdd             Strings     `yaml:"cap_add,omitempty"`
	CapDrop            Strings     `yaml:"cap_drop,omitempty"`
	CgroupParent       String      `yaml:"cgroup_parent,omitempty"`
	Command            Command     `yaml:"command,omitempty"`
	Configs            Strings     `yaml:"configs,omitempty"` // TODO: support long syntax.
	ContainerName      String      `yaml:"container_name,omitempty"`
	// TODO: credential_spec
	DependsOn         ServiceDependencies     `yaml:"depends_on,omitempty"`
	DeviceCgroupRules Strings                 `yaml:"device_cgroup_rules,omitempty"`
	Devices           []DeviceMapping         `yaml:"devices,omitempty"`
	DNS               Tuple                   `yaml:"dns,omitempty"`
	DNSOptions        Strings                 `yaml:"dns_opt,omitempty"`
	DNSSearch         Tuple                   `yaml:"dns_search,omitempty"`
	Domainname        String                  `yaml:"domainname,omitempty"`
	Entrypoint        Command                 `yaml:"entrypoint,omitempty"`
	EnvFile           Tuple                   `yaml:"env_file,omitempty"`
	Environment       Dictionary              `yaml:"environment,omitempty"`
	Expose            []PortRangeWithProtocol `yaml:"expose,omitempty"`
	// TODO: extends
	// List of links of the form `SERVICE` or `SERVICE:ALIAS`
	ExternalLinks Strings `yaml:"external_links,omitempty"`
	// List of host/IP pairs to add to /etc/hosts of the form `HOST:IP`
	ExtraHosts       Strings         `yaml:"extra_hosts,omitempty"`
	GroupAdd         Strings         `yaml:"group_add,omitempty"`
	Healthcheck      *Healthcheck    `yaml:"healthcheck,omitempty"`
	Hostname         String          `yaml:"hostname,omitempty"`
	Image            String          `yaml:"image,omitempty"`
	Init             *Bool           `yaml:"init,omitempty"`
	IPC              String          `yaml:"ipc,omitempty"`
	Isolation        String          `yaml:"isolation,omitempty"`
	Labels           Dictionary      `yaml:"labels,omitempty"`
	Links            Links           `yaml:"links,omitempty"`
	Logging          Logging         `yaml:"logging,omitempty"`
	NetworkMode      String          `yaml:"network_mode,omitempty"`
	Networks         ServiceNetworks `yaml:"networks,omitempty"`
	MacAddress       String          `yaml:"mac_address,omitempty"`
	MemorySwappiness *Int            `yaml:"mem_swappiness,omitempty"`
	// MemoryLimit and MemoryReservation can be specified either as strings or integers.
	// TODO: Deprecate these fields once we support `deploy.limits.memory` and `deploy.reservations.memory`.
	MemoryLimit       Bytes `yaml:"mem_limit,omitempty"`
	MemoryReservation Bytes `yaml:"mem_reservation,omitempty"`

	MemswapLimit   Bytes        `yaml:"memswap_limit,omitempty"`
	OomKillDisable *Bool        `yaml:"oom_kill_disable,omitempty"`
	OomScoreAdj    Int          `yaml:"oom_score_adj,omitempty"`
	PidMode        String       `yaml:"pid,omitempty"`
	PidsLimit      *Int         `yaml:"pids_limit,omitempty"`
	Platform       String       `yaml:"platform,omitempty"`
	Ports          PortMappings `yaml:"ports,omitempty"`
	Privileged     Bool         `yaml:"privileged,omitempty"`
	// TODO: Support profiles. See https://docs.docker.com/compose/profiles/.
	Profiles        Ignored       `yaml:"profiles,omitempty"`
	PullPolicy      String        `yaml:"pull_policy,omitempty"`
	ReadOnly        Bool          `yaml:"read_only,omitempty"`
	Restart         String        `yaml:"restart,omitempty"`
	Runtime         String        `yaml:"runtime,omitempty"`
	SecurityOpt     Strings       `yaml:"security_opt,omitempty"`
	ShmSize         Bytes         `yaml:"shm_size,omitempty"`
	StdinOpen       Bool          `yaml:"stdin_open,omitempty"`
	StopGracePeriod *Duration     `yaml:"stop_grace_period,omitempty"`
	StopSignal      String        `yaml:"stop_signal,omitempty"`
	StorageOpt      Dictionary    `yaml:"storage_opt,omitempty"`
	Sysctls         Dictionary    `yaml:"sysctls,omitempty"`
	Tmpfs           Tuple         `yaml:"tmpfs,omitempty"`
	TTY             Bool          `yaml:"tty,omitempty"`
	Ulimits         Ulimits       `yaml:"ulimits,omitempty"`
	User            String        `yaml:"user,omitempty"`
	UsernsMode      String        `yaml:"userns_mode,omitempty"`
	Volumes         []VolumeMount `yaml:"volumes,omitempty"`
	VolumesFrom     Strings       `yaml:"volumes_from,omitempty"`
	WorkingDir      String        `yaml:"working_dir,omitempty"`

	// NOTE [DOCKER SWARM FEATURES]:
	// Docker-Compose manages local, single-container deployments as well as Docker Swarm
	// deployments. Since Swarm is not as widely used as Kubernetes, support for the Swarm
	// features that Docker-Compose includes is not a top priority. The settings listed
	// below are the ones that are applicable to a Swarm deployment.
	Deploy  Ignored `yaml:"deploy,omitempty"`
	Scale   Ignored `yaml:"scale,omitempty"`
	Secrets Ignored `yaml:"secrets,omitempty"`
}

func (service *Service) Interpolate(env Environment) error {
	return interpolateStruct(service, env)
}
