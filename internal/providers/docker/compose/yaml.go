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
	"regexp"
	"strings"

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
	Volumes []VolumeMount `yaml:"volumes"`
	// TODO: volumes_from

	WorkingDir string `yaml:"working_dir"`
}

type VolumeMount struct {
	Type        string
	Source      string
	Target      string
	ReadOnly    bool
	Bind        BindOptions
	Volume      VolumeOptions
	Tmpfs       TmpfsOptions
	Consistency IgnoredField
}

func (vm *VolumeMount) UnmarshalYAML(b []byte) error {
	var asString string
	if err := yaml.Unmarshal(b, &asString); err == nil {
		return vm.fromShortSyntax(asString)
	}

	asExtended := extendedVolumeMount{}
	if err := yaml.Unmarshal(b, &asExtended); err != nil {
		return err
	}
	vm.Type = asExtended.Type
	vm.Source = asExtended.Source
	vm.Target = asExtended.Target
	vm.ReadOnly = asExtended.ReadOnly
	vm.Bind = asExtended.Bind
	vm.Volume = asExtended.Volume
	vm.Tmpfs = asExtended.Tmpfs
	vm.Consistency = asExtended.Consistency

	return nil
}

func (vm *VolumeMount) fromShortSyntax(in string) error {
	parts := strings.Split(in, ":")
	switch len(parts) {
	case 1:
		vm.Type = "volume"
		vm.Target = in
	case 2:
		vm.setSource(parts[0])
		vm.Target = parts[1]
	case 3:
		vm.setSource(parts[0])
		vm.Target = parts[1]
		accessMode := parts[2]
		switch accessMode {
		case "ro":
			vm.ReadOnly = true
		case "rw":
			// Do nothing - va.ReadOnly is already false.
		default:
			return fmt.Errorf(`invalid access mode; expected "ro" or "rw" but got %q`, accessMode)
		}
	default:
		return fmt.Errorf(`invalid volume specification; expected "VOLUME:CONTAINER_PATH" or "VOLUME:CONTAINER_PATH:ACCESS_MODE" but got %q`, in)
	}

	return nil
}

var localPathRe = regexp.MustCompile("^[./~]")

func (vm *VolumeMount) setSource(src string) {
	vm.Source = src
	if localPathRe.MatchString(src) {
		vm.Type = "bind"
		vm.Bind = BindOptions{
			// CreateHostPath is always implied by the short syntax.
			CreateHostPath: true,
		}
	} else {
		vm.Type = "volume"
	}
}

// extendedVolumeMount is a private struct that is structurally identical to VolumeAttachment but
// is only used for YAML unmarshalling where we do not need to consider the short string-based syntax.
type extendedVolumeMount struct {
	Type        string        `yaml:"type"`
	Source      string        `yaml:"source"`
	Target      string        `yaml:"target"`
	ReadOnly    bool          `yaml:"read_only"`
	Bind        BindOptions   `yaml:"bind"`
	Volume      VolumeOptions `yaml:"volume"`
	Tmpfs       TmpfsOptions  `yaml:"tmpfs"`
	Consistency IgnoredField  `yaml:"consistency"`
}

type VolumeOptions struct {
	Nocopy bool `yaml:"nocopy"`
}

type BindOptions struct {
	Propagation    string `yaml:"propagation"`
	CreateHostPath bool   `yaml:"create_host_path"`
}

type TmpfsOptions struct {
	Size int64 `yaml:"size"`
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
