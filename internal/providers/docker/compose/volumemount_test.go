package compose

import "testing"

func TestVolumeMountYAML(t *testing.T) {
	testYAML(t, "short_volume", "target", VolumeMount{
		ShortForm: MakeString("target"),
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   MakeString("volume"),
			Target: MakeString("target"),
		},
	})
	testYAML(t, "short_bind", "volume:/container_path", VolumeMount{
		ShortForm: MakeString("volume:/container_path"),
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   MakeString("volume"),
			Source: MakeString("volume"),
			Target: MakeString("/container_path"),
		},
	})
	testYAML(t, "short_access", "./host_path:/container_path:ro", VolumeMount{
		ShortForm: MakeString("./host_path:/container_path:ro"),
		VolumeMountLongForm: VolumeMountLongForm{
			Type:     MakeString("bind"),
			Source:   MakeString("./host_path"),
			Target:   MakeString("/container_path"),
			ReadOnly: MakeBool(true),
			Bind: &BindOptions{
				CreateHostPath: MakeBool(true),
			},
		},
	})
	testYAML(t, "long_1", `
type: volume
source: mydata
target: /data
read_only: true
volume:
  nocopy: true
`, VolumeMount{
		VolumeMountLongForm: VolumeMountLongForm{
			Type:     MakeString("volume"),
			Source:   MakeString("mydata"),
			Target:   MakeString("/data"),
			ReadOnly: MakeBool(true),
			Volume: &VolumeOptions{
				Nocopy: MakeBool(true),
			},
		},
	})
	testYAML(t, "long_2", `
type: bind
source: /path/a
target: /path/b
bind:
  propagation: rshared
  create_host_path: true
`, VolumeMount{
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   MakeString("bind"),
			Source: MakeString("/path/a"),
			Target: MakeString("/path/b"),
			Bind: &BindOptions{
				Propagation:    MakeString("rshared"),
				CreateHostPath: MakeBool(true),
			},
		},
	})
	testYAML(t, "long_3", `
type: tmpfs
target: /data/buffer
tmpfs:
  size: 208666624
`, VolumeMount{
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   MakeString("tmpfs"),
			Target: MakeString("/data/buffer"),
			Tmpfs: &TmpfsOptions{
				Size: Bytes{
					String:   MakeInt(208666624).String,
					Quantity: 208666624,
				},
			},
		},
	})

	assertInterpolated(t, map[string]string{
		"short": "./host_path:/container_path:ro",
	}, "${short}", VolumeMount{
		ShortForm: MakeString("${short}").WithValue("./host_path:/container_path:ro"),
		VolumeMountLongForm: VolumeMountLongForm{
			Type:     MakeString("bind"),
			Source:   MakeString("./host_path"),
			Target:   MakeString("/container_path"),
			ReadOnly: MakeBool(true),
			Bind: &BindOptions{
				CreateHostPath: MakeBool(true),
			},
		},
	})
}
