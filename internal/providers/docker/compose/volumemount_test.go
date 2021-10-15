package compose

import "testing"

func TestVolumeMountYAML(t *testing.T) {
	testYAML(t, "short_volume", "target", VolumeMount{
		IsShortForm: true,
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   "volume",
			Target: "target",
		},
	})
	testYAML(t, "short_bind", "volume:/container_path", VolumeMount{
		IsShortForm: true,
		VolumeMountLongForm: VolumeMountLongForm{
			Type:   "volume",
			Source: "volume",
			Target: "/container_path",
		},
	})
	testYAML(t, "short_access", "./host_path:/container_path:ro", VolumeMount{
		IsShortForm: true,
		VolumeMountLongForm: VolumeMountLongForm{
			Type:     "bind",
			Source:   "./host_path",
			Target:   "/container_path",
			ReadOnly: true,
			Bind: &BindOptions{
				CreateHostPath: true,
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
			Type:     "volume",
			Source:   "mydata",
			Target:   "/data",
			ReadOnly: true,
			Volume: &VolumeOptions{
				Nocopy: true,
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
			Type:   "bind",
			Source: "/path/a",
			Target: "/path/b",
			Bind: &BindOptions{
				Propagation:    "rshared",
				CreateHostPath: true,
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
			Type:   "tmpfs",
			Target: "/data/buffer",
			Tmpfs: &TmpfsOptions{
				Size: Bytes{
					Quantity: 208666624,
				},
			},
		},
	})
}
