package compose

import "testing"

func TestServiceNetworksYAML(t *testing.T) {
	testYAML(t, "names", `
- one
- two
`, ServiceNetworks{
		Style: SeqStyle,
		Items: []ServiceNetwork{
			ServiceNetwork{
				Name:        "one",
				IsShortForm: true,
			},
			ServiceNetwork{
				Name:        "two",
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
			ServiceNetwork{
				Name: "one",
			},
			ServiceNetwork{
				Name: "two",
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
