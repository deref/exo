package compose

import "testing"

func TestProjectYAML(t *testing.T) {
	testYAML(t, "empty", `{}`, Project{})
	testYAML(t, "sections", `
version: 3
services:
  one:
    image: service
  two: {}
networks:
  three:
    name: network
  four: {}
volumes:
  five:
    name: volume
  six: {}
configs:
  seven:
    name: config
  eight: {}
secrets:
  nine:
    name: secret
  ten: {}
`, Project{
		Version: String(MakeInt(3)),
		Services: ProjectServices{
			{
				Key:   "one",
				Image: "service",
			},
			{
				Key: "two",
			},
		},
		Networks: ProjectNetworks{
			{
				Key:  "three",
				Name: "network",
			},
			{
				Key: "four",
			},
		},
		Volumes: ProjectVolumes{
			{
				Key:  "five",
				Name: "volume",
			},
			{
				Key: "six",
			},
		},
		Configs: ProjectConfigs{
			{
				Key:  "seven",
				Name: "config",
			},
			{
				Key: "eight",
			},
		},
		Secrets: ProjectSecrets{
			{
				Key:  "nine",
				Name: "secret",
			},
			{
				Key: "ten",
			},
		},
	})
}
