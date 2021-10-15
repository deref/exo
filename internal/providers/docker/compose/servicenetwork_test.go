package compose

import "testing"

func TestServiceNetworksYAML(t *testing.T) {
	testYAML(t, "short", `net`, ServiceNetwork{
		Key:         "net",
		IsShortForm: true,
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
			Aliases:      []string{"foo", "bar"},
			IPV4Address:  "172.16.238.10",
			IPV6Address:  "2001:3984:3989::10",
			LinkLocalIPs: []string{"57.123.22.11", "57.123.22.13"},
			Priority:     1000,
		},
	})

	testYAML(t, "seq", `
- one
- two
`, ServiceNetworks{
		Style: SeqStyle,
		Items: []ServiceNetwork{
			{
				Key:         "one",
				IsShortForm: true,
			},
			{
				Key:         "two",
				IsShortForm: true,
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
					Aliases: []string{
						"a",
						"b",
					},
				},
			},
		},
	})
}
