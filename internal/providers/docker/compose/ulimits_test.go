package compose

import "testing"

func TestUlimitYAML(t *testing.T) {
	testYAML(t, "short", `1024`, Ulimit{
		ShortForm: MakeInt(1024),
		UlimitLongForm: UlimitLongForm{
			Soft: MakeInt(1024),
			Hard: MakeInt(1024),
		},
	})
	testYAML(t, "long", `
soft: 1024
hard: 2048
`, Ulimit{
		UlimitLongForm: UlimitLongForm{
			Soft: MakeInt(1024),
			Hard: MakeInt(2048),
		},
	})
	testYAML(t, "map", `
nproc: 65535
nofile:
  soft: 20000
  hard: 40000
`, Ulimits{
		{
			Name:      "nproc",
			ShortForm: MakeInt(65535),
			UlimitLongForm: UlimitLongForm{
				Soft: MakeInt(65535),
				Hard: MakeInt(65535),
			},
		},
		{
			Name: "nofile",
			UlimitLongForm: UlimitLongForm{
				Soft: MakeInt(20000),
				Hard: MakeInt(40000),
			},
		},
	})
	assertInterpolated(t, map[string]string{
		"x":    "2048",
		"soft": "4096",
	}, `
nproc: ${x}
nofile:
  soft: ${soft}
  hard: 40000
	`, Ulimits{
		{
			Name: "nproc",
			ShortForm: Int{
				String: MakeString("${x}").WithValue("2048"),
				Value:  2048,
			},
			UlimitLongForm: UlimitLongForm{
				Soft: Int{
					String: MakeString("${x}").WithValue("2048"),
					Value:  2048,
				},
				Hard: Int{
					String: MakeString("${x}").WithValue("2048"),
					Value:  2048,
				},
			},
		},
		{
			Name: "nofile",
			UlimitLongForm: UlimitLongForm{
				Soft: Int{
					String: MakeString("${soft}").WithValue("4096"),
					Value:  4096,
				},
				Hard: MakeInt(40000),
			},
		},
	})
}
