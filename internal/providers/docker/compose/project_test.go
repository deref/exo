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
		Version: MakeInt(3).String,
		Services: ProjectServices{
			{
				Key:   "one",
				Image: MakeString("service"),
			},
			{
				Key: "two",
			},
		},
		Networks: ProjectNetworks{
			{
				Key:  "three",
				Name: MakeString("network"),
			},
			{
				Key: "four",
			},
		},
		Volumes: ProjectVolumes{
			{
				Key:  "five",
				Name: MakeString("volume"),
			},
			{
				Key: "six",
			},
		},
		Configs: ProjectConfigs{
			{
				Key:  "seven",
				Name: MakeString("config"),
			},
			{
				Key: "eight",
			},
		},
		Secrets: ProjectSecrets{
			{
				Key:  "nine",
				Name: MakeString("secret"),
			},
			{
				Key: "ten",
			},
		},
	})

	assertInterpolated(t, map[string]string{
		"service": "SERVICE",
		"network": "NETWORK",
		"volume":  "VOLUME",
		"config":  "CONFIG",
		"secret":  "SECRET",
	}, `
version: 3
services:
  one:
    image: ${service}
  two: {}
networks:
  three:
    name: ${network}
  four: {}
volumes:
  five:
    name: ${volume}
  six: {}
configs:
  seven:
    name: ${config}
  eight: {}
secrets:
  nine:
    name: ${secret}
  ten: {}
`, Project{
		Version: MakeInt(3).String,
		Services: ProjectServices{
			{
				Key:   "one",
				Image: MakeString("${service}").WithValue("SERVICE"),
			},
			{
				Key: "two",
			},
		},
		Networks: ProjectNetworks{
			{
				Key:  "three",
				Name: MakeString("${network}").WithValue("NETWORK"),
			},
			{
				Key: "four",
			},
		},
		Volumes: ProjectVolumes{
			{
				Key:  "five",
				Name: MakeString("${volume}").WithValue("VOLUME"),
			},
			{
				Key: "six",
			},
		},
		Configs: ProjectConfigs{
			{
				Key:  "seven",
				Name: MakeString("${config}").WithValue("CONFIG"),
			},
			{
				Key: "eight",
			},
		},
		Secrets: ProjectSecrets{
			{
				Key:  "nine",
				Name: MakeString("${secret}").WithValue("SECRET"),
			},
			{
				Key: "ten",
			},
		},
	})
}
