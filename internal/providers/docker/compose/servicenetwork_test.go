package compose

import "testing"

func TestServiceNetworksYAML(t *testing.T) {
	testYAML(t, "short", `net`, ServiceNetwork{
		Key:       "net",
		ShortForm: MakeString("net"),
	})
	testYAML(t, "long", `
aliases:
  - foo
  - bar
ipv4_address: 172.16.238.10
ipv6_address: 2001:3984:3989::10
link_local_ips:
  - 57.123.22.11
  - 57.123.22.13
priority: 1000
	`, ServiceNetwork{
		ServiceNetworkLongForm: ServiceNetworkLongForm{
			Aliases:     []String{MakeString("foo"), MakeString("bar")},
			IPV4Address: MakeString("172.16.238.10"),
			IPV6Address: MakeString("2001:3984:3989::10"),
			LinkLocalIPs: []String{
				MakeString("57.123.22.11"),
				MakeString("57.123.22.13"),
			},
			Priority: MakeInt(1000),
		},
	})

	testYAML(t, "seq", `
- one
- two
`, ServiceNetworks{
		Style: SeqStyle,
		Items: []ServiceNetwork{
			{
				Key:       "one",
				ShortForm: MakeString("one"),
			},
			{
				Key:       "two",
				ShortForm: MakeString("two"),
			},
		},
	})
	testYAML(t, "map", `
one: {}
two:
  aliases:
    - a
    - b
`, ServiceNetworks{
		Style: MapStyle,
		Items: []ServiceNetwork{
			{
				Key: "one",
			},
			{
				Key: "two",
				ServiceNetworkLongForm: ServiceNetworkLongForm{
					Aliases: []String{
						MakeString("a"),
						MakeString("b"),
					},
				},
			},
		},
	})
	assertInterpolated(t, map[string]string{"network": "NETWORK"}, `
- ${network}
`, ServiceNetworks{
		Style: SeqStyle,
		Items: []ServiceNetwork{
			{
				Key:       "NETWORK",
				ShortForm: MakeString("${network}").WithValue("NETWORK"),
			},
		},
	})
	assertInterpolated(t, map[string]string{"alias": "ALIAS"}, `
key:
  aliases:
    - ${alias}
`, ServiceNetworks{
		Style: MapStyle,
		Items: []ServiceNetwork{
			{
				Key: "key",
				ServiceNetworkLongForm: ServiceNetworkLongForm{
					Aliases: Strings{
						MakeString("${alias}").WithValue("ALIAS"),
					},
				},
			},
		},
	})
}
