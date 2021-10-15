package compose

import (
	"testing"
	"time"
)

func TestServiceYAML(t *testing.T) {
	testYAML(t, "cpu_config", `
cpu_count: 2
cpu_percent: 80
cpu_rt_runtime: 400ms
cpu_rt_period: 1400
cpuset: 0,2,4
`, Service{
		CPUCount:   2,
		CPUPercent: 80,
		CPURealtimeRuntime: Duration{
			Expression: "400ms",
			Duration:   400 * time.Millisecond,
		},
		CPURealtimePeriod: Duration{
			Expression: "1400",
			Duration:   1400 * time.Microsecond,
		},
		CPUSet: "0,2,4",
	})

	testYAML(t, "capabilities", `
cap_add:
  - ALL
cap_drop:
  - NET_ADMIN
  - SYS_ADMIN
`, Service{
		CapAdd:  []string{"ALL"},
		CapDrop: []string{"NET_ADMIN", "SYS_ADMIN"},
	})

	testYAML(t, "cgroup_parent", `
cgroup_parent: m-executor-abcd
`, Service{
		CgroupParent: "m-executor-abcd",
	})

	testYAML(t, "cgroup_rules", `
device_cgroup_rules:
  - c 1:3 mr
  - a 7:* rmw
`, Service{
		DeviceCgroupRules: []string{
			"c 1:3 mr",
			"a 7:* rmw",
		},
	})

	testYAML(t, "dns_single", `
dns: 8.8.8.8
`, Service{
		DNS: MakeStrings("8.8.8.8"),
	})

	testYAML(t, "dns_multiple", `
dns:
  - 8.8.8.8
  - 4.4.4.4
`, Service{
		DNS: MakeStrings("8.8.8.8", "4.4.4.4"),
	})

	testYAML(t, "dns_options", `
dns_opt:
  - use-vc
  - no-tld-query
`, Service{
		DNSOptions: []string{"use-vc", "no-tld-query"},
	})

	testYAML(t, "dns_search_short", `
dns_search: example.com
`, Service{
		DNSSearch: MakeStrings("example.com"),
	})

	testYAML(t, "dns_search_long", `
dns_search:
  - ns1.example.com
  - ns2.example.com
`, Service{
		DNSSearch: MakeStrings("ns1.example.com", "ns2.example.com"),
	})

	testYAML(t, "env_file", `
env_file: .dockerenv
`, Service{
		EnvFile: MakeStrings(".dockerenv"),
	})

	testYAML(t, "external_links", `
external_links:
  - container1
  - container2:alias
`, Service{
		ExternalLinks: []string{"container1", "container2:alias"},
	})

	testYAML(t, "extra_hosts", `
extra_hosts:
  - somehost:162.242.195.82
  - otherhost:50.31.209.229
`, Service{
		ExtraHosts: []string{"somehost:162.242.195.82", "otherhost:50.31.209.229"},
	})

	testYAML(t, "group_add", `
group_add:
  - mail
`, Service{
		GroupAdd: []string{"mail"},
	})

	testYAML(t, "init", `
init: true
`, Service{
		Init: boolRef(true),
	})

	testYAML(t, "ipc", `
ipc: service:foo
`, Service{
		IPC: "service:foo",
	})

	testYAML(t, "isolation", `
isolation: hyperv
`, Service{
		Isolation: "hyperv",
	})

	testYAML(t, "network_mode", `
network_mode: host
`, Service{
		NetworkMode: "host",
	})

	testYAML(t, "memswap_limit", `
memswap_limit: 2g
`,
		Service{
			MemswapLimit: Bytes{
				Quantity: 2,
				Unit: ByteUnit{
					Suffix: "g",
					Scalar: 1024 * 1024 * 1024,
				},
			},
		})

	testYAML(t, "oom_kill_disable", `
oom_kill_disable: true
oom_score_adj: 200
`, Service{
		OomKillDisable: boolRef(true),
		OomScoreAdj:    200,
	})

	testYAML(t, "pid", `
pid: host
pids_limit: 5280
`, Service{
		PidMode:   "host",
		PidsLimit: int64Ref(5280),
	})

	testYAML(t, "platform", `
platform: linux/arm64/v8
`, Service{
		Platform: "linux/arm64/v8",
	})

	testYAML(t, "pull_policy", `
pull_policy: missing
`, Service{
		PullPolicy: "missing",
	})

	testYAML(t, "read_only", `
read_only: true
`, Service{
		ReadOnly: true,
	})

	testYAML(t, "storage_opt", `
storage_opt:
  size: 20G
`, Service{
		StorageOpt: Dictionary{
			Style: MapStyle,
			Items: []DictionaryItem{
				{
					Style: MapStyle,
					Key:   "size",
					Value: "20G",
				},
			},
		},
	})

	testYAML(t, "sysctls", `
sysctls:
  - net.core.somaxconn=1024
  - net.ipv4.tcp_syncookies=0
`, Service{
		Sysctls: Dictionary{
			Style: SeqStyle,
			Items: []DictionaryItem{
				{
					Style: SeqStyle,
					Key:   "net.core.somaxconn",
					Value: "1024",
				},
				{
					Style: SeqStyle,
					Key:   "net.ipv4.tcp_syncookies",
					Value: "0",
				},
			},
		},
	})

	testYAML(t, "userns_mode", `
userns_mode: host
`,
		Service{
			UsernsMode: "host",
		})

	testYAML(t, "volumes_from", `
volumes_from:
  - container:my-container:ro
`, Service{
		VolumesFrom: []string{"container:my-container:ro"},
	})

}
