package compose_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/stretchr/testify/assert"
)

func TestParseService(t *testing.T) {
	trueVal := true
	int64Val := int64(5280)
	testCases := []struct {
		name     string
		in       string
		expected compose.Service
	}{
		{
			name: "volumes - full syntax",
			in: `volumes:
- type: volume
  source: mydata
  target: /data
  read_only: true
  volume:
    nocopy: true
- type: bind
  source: /path/a
  target: /path/b
  bind:
    propagation: rshared
    create_host_path: true
- type: tmpfs
  target: /data/buffer
  tmpfs:
    size: 208666624`,
			expected: compose.Service{
				Volumes: []compose.VolumeMount{
					{
						Type:     "volume",
						Source:   "mydata",
						Target:   "/data",
						ReadOnly: true,
						Volume: &compose.VolumeOptions{
							Nocopy: true,
						},
					},
					{
						Type:   "bind",
						Source: "/path/a",
						Target: "/path/b",
						Bind: &compose.BindOptions{
							Propagation:    "rshared",
							CreateHostPath: true,
						},
					},
					{
						Type:   "tmpfs",
						Target: "/data/buffer",
						Tmpfs: &compose.TmpfsOptions{
							Size: 208666624,
						},
					},
				},
			},
		},

		{
			name: "volumes: short syntax",
			in: `volumes:
- /var/myapp
- './data:/data'
- "/home/fred/.ssh:/root/.ssh:ro"
- '~/util:/usr/bin/util:rw'
- my-log-volume:/var/log/xyzzy`,
			expected: compose.Service{
				Volumes: []compose.VolumeMount{
					{
						Type:   "volume",
						Target: "/var/myapp",
					},
					{
						Type:   "bind",
						Source: "./data",
						Target: "/data",
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:     "bind",
						Source:   "/home/fred/.ssh",
						Target:   "/root/.ssh",
						ReadOnly: true,
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:   "bind",
						Source: "~/util",
						Target: "/usr/bin/util",
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:   "volume",
						Source: "my-log-volume",
						Target: "/var/log/xyzzy",
					},
				},
			},
		},

		{
			name: "service dependencies - short syntax",
			in: `depends_on:
- db
- messages`,
			expected: compose.Service{
				DependsOn: compose.ServiceDependencies{
					IsShortSyntax: true,
					Services: []compose.ServiceDependency{
						{
							Service:   "db",
							Condition: "service_started",
						},
						{
							Service:   "messages",
							Condition: "service_started",
						},
					},
				},
			},
		},

		{
			name: "service dependencies - extended syntax",
			in: `depends_on:
  db:
  messages:
    condition: service_healthy`,
			expected: compose.Service{
				DependsOn: compose.ServiceDependencies{
					Services: []compose.ServiceDependency{
						{
							Service:   "db",
							Condition: "service_started",
						},
						{
							Service:   "messages",
							Condition: "service_healthy",
						},
					},
				},
			},
		},

		{
			name: "cpu config",
			in: `cpu_count: 2
cpu_percent: 80
cpu_rt_runtime: 400ms
cpu_rt_period: 1400
cpuset: 0,2,4`,
			expected: compose.Service{
				CPUCount:           2,
				CPUPercent:         80,
				CPURealtimeRuntime: compose.Duration(400 * time.Millisecond),
				CPURealtimePeriod:  compose.Duration(1400 * time.Microsecond),
				CPUSet:             "0,2,4",
			},
		},

		{
			name: "block io config",
			in: `blkio_config:
  weight: 300
  weight_device:
  - path: /dev/sda
    weight: 400
  device_read_bps:
  - path: /dev/sdb
    rate: '12mb'
  device_read_iops:
  - path: /dev/sdb
    rate: 120
  device_write_bps:
  - path: /dev/sdb
    rate: '1024k'
  device_write_iops:
  - path: /dev/sdb
    rate: 30`,
			expected: compose.Service{
				BlkioConfig: compose.BlkioConfig{
					Weight: 300,
					WeightDevice: []compose.WeightDevice{
						{
							Path:   "/dev/sda",
							Weight: 400,
						},
					},
					DeviceReadBPS: []compose.ThrottleDevice{
						{
							Path: "/dev/sdb",
							Rate: compose.Bytes(12582912),
						},
					},
					DeviceReadIOPS: []compose.ThrottleDevice{
						{
							Path: "/dev/sdb",
							Rate: 120,
						},
					},
					DeviceWriteBPS: []compose.ThrottleDevice{
						{
							Path: "/dev/sdb",
							Rate: compose.Bytes(1048576),
						},
					},
					DeviceWriteIOPS: []compose.ThrottleDevice{
						{
							Path: "/dev/sdb",
							Rate: 30,
						},
					},
				},
			},
		},

		{
			name: "capabilities",
			in: `cap_add:
- ALL
cap_drop:
- NET_ADMIN
- SYS_ADMIN`,
			expected: compose.Service{
				CapAdd:  []string{"ALL"},
				CapDrop: []string{"NET_ADMIN", "SYS_ADMIN"},
			},
		},

		{
			name: "cgroup parent",
			in:   `cgroup_parent: m-executor-abcd`,
			expected: compose.Service{
				CgroupParent: "m-executor-abcd",
			},
		},

		{
			name: "device cgroup rules",
			in: `device_cgroup_rules:
- 'c 1:3 mr'
- 'a 7:* rmw'`,
			expected: compose.Service{
				DeviceCgroupRules: []string{"c 1:3 mr", "a 7:* rmw"},
			},
		},

		{
			name: "device mappings",
			in: `devices:
- "/dev/ttyUSB0:/dev/ttyUSB1"`,
			expected: compose.Service{
				Devices: []compose.DeviceMapping{
					{
						PathOnHost:      "/dev/ttyUSB0",
						PathInContainer: "/dev/ttyUSB1",
					},
				},
			},
		},

		{
			name: "devices",
			in: `devices:
- "/dev/ttyUSB0:/dev/ttyUSB1"
- "/dev/sda:/dev/xvda:rwm"`,
			expected: compose.Service{
				Devices: []compose.DeviceMapping{
					{
						PathOnHost:      "/dev/ttyUSB0",
						PathInContainer: "/dev/ttyUSB1",
					},
					{
						PathOnHost:        "/dev/sda",
						PathInContainer:   "/dev/xvda",
						CgroupPermissions: "rwm",
					},
				},
			},
		},

		{
			name: "DNS - single string",
			in:   `dns: "8.8.8.8"`,
			expected: compose.Service{
				DNS: compose.StringOrStringSlice{"8.8.8.8"},
			},
		},

		{
			name: "DNS - list",
			in: `dns:
- '8.8.8.8'
- '4.4.4.4'`,
			expected: compose.Service{
				DNS: compose.StringOrStringSlice{"8.8.8.8", "4.4.4.4"},
			},
		},

		{
			name: "DNS options",
			in: `dns_opt:
- use-vc
- no-tld-query`,
			expected: compose.Service{
				DNSOptions: []string{"use-vc", "no-tld-query"},
			},
		},

		{
			name: "DNS search - short",
			in:   "dns_search: example.com",
			expected: compose.Service{
				DNSSearch: compose.StringOrStringSlice{"example.com"},
			},
		},

		{
			name: "DNS search - long",
			in: `dns_search:
- ns1.example.com
- ns2.example.com`,
			expected: compose.Service{
				DNSSearch: compose.StringOrStringSlice{"ns1.example.com", "ns2.example.com"},
			},
		},

		{
			name: "Env files",
			in:   `env_file: .dockerenv`,
			expected: compose.Service{
				EnvFile: compose.StringOrStringSlice{".dockerenv"},
			},
		},

		{
			name: "External links",
			in: `external_links:
- container1
- container2:alias`,
			expected: compose.Service{
				ExternalLinks: []string{"container1", "container2:alias"},
			},
		},

		{
			name: "Extra hosts",
			in: `extra_hosts:
- "somehost:162.242.195.82"
- "otherhost:50.31.209.229"`,
			expected: compose.Service{
				ExtraHosts: []string{"somehost:162.242.195.82", "otherhost:50.31.209.229"},
			},
		},

		{
			name: "Additional groups",
			in: `group_add:
- mail`,
			expected: compose.Service{
				GroupAdd: []string{"mail"},
			},
		},

		{
			name: "Init process flag",
			in:   `init: true`,
			expected: compose.Service{
				Init: &trueVal,
			},
		},

		{
			name: "IPC",
			in:   `ipc: "service:foo"`,
			expected: compose.Service{
				IPC: "service:foo",
			},
		},

		{
			name: "Isolation",
			in:   `isolation: hyperv`,
			expected: compose.Service{
				Isolation: "hyperv",
			},
		},

		{
			name: "Network mode",
			in:   `network_mode: host`,
			expected: compose.Service{
				NetworkMode: "host",
			},
		},

		{
			name: "Networks - short form",
			in: `networks:
- net1
- net2`,
			expected: compose.Service{
				Networks: compose.ServiceNetworks{
					{
						Network: "net1",
					},
					{
						Network: "net2",
					},
				},
			},
		},

		{
			name: "Networks - long form",
			in: `networks:
  net1:
  net2:
    aliases:
    - foo
    - bar
    ipv4_address: 172.16.238.10
    ipv6_address: 2001:3984:3989::10
    link_local_ips:
    - 57.123.22.11
    - 57.123.22.13
    priority: 1000`,
			expected: compose.Service{
				Networks: compose.ServiceNetworks{
					{
						Network: "net1",
					},
					{
						Network:      "net2",
						Aliases:      []string{"foo", "bar"},
						IPV4Address:  "172.16.238.10",
						IPV6Address:  "2001:3984:3989::10",
						LinkLocalIPs: []string{"57.123.22.11", "57.123.22.13"},
						Priority:     1000,
					},
				},
			},
		},

		{
			name: "Memswap Limit - numeric",
			in:   `memswap_limit: -1`,
			expected: compose.Service{
				MemswapLimit: -1,
			},
		},

		{
			name: "Memswap Limit - bytes string",
			in:   `memswap_limit: 2g`,
			expected: compose.Service{
				MemswapLimit: 2147483648,
			},
		},

		{
			name: "OOM Settings",
			in: `oom_kill_disable: true
oom_score_adj: 200`,
			expected: compose.Service{
				OomKillDisable: &trueVal,
				OomScoreAdj:    200,
			},
		},

		{
			name: "PID settings",
			in: `pid: host
pids_limit: 5280`,
			expected: compose.Service{
				PidMode:   "host",
				PidsLimit: &int64Val,
			},
		},

		{
			name: "Platform",
			in:   `platform: linux/arm64/v8`,
			expected: compose.Service{
				Platform: "linux/arm64/v8",
			},
		},

		{
			name: "Pull Policy",
			in:   `pull_policy: missing`,
			expected: compose.Service{
				PullPolicy: "missing",
			},
		},
	}

	for _, testCase := range testCases {
		name := testCase.name
		in := testCase.in
		expected := testCase.expected
		t.Run(name, func(t *testing.T) {
			var content bytes.Buffer
			content.WriteString("services:\n  test-svc:\n")
			lines := strings.Split(in, "\n")
			for _, line := range lines {
				content.WriteString("    ")
				content.WriteString(line)
				content.WriteByte('\n')
			}

			comp, err := compose.Parse(&content)
			if !assert.NoError(t, err) {
				return
			}

			svc := comp.Services["test-svc"]
			assert.Equal(t, expected, svc)
		})
	}
}
