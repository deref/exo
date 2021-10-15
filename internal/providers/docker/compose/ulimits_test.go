package compose

import "testing"

func TestUlimitYAML(t *testing.T) {
	testYAML(t, "short", `1024`, Ulimit{
		IsShortForm: true,
		UlimitLongForm: UlimitLongForm{
			Soft: 1024,
			Hard: 1024,
		},
	})
	testYAML(t, "long", `
soft: 1024
hard: 2048
`, Ulimit{
		UlimitLongForm: UlimitLongForm{
			Soft: 1024,
			Hard: 2048,
		},
	})
	testYAML(t, "map", `
nproc: 65535
nofile:
  soft: 20000
  hard: 40000
`, Ulimits{
		{
			Name:        "nproc",
			IsShortForm: true,
			UlimitLongForm: UlimitLongForm{
				Soft: 65535,
				Hard: 65535,
			},
		},
		{
			Name: "nofile",
			UlimitLongForm: UlimitLongForm{
				Soft: 20000,
				Hard: 40000,
			},
		},
	})
}
