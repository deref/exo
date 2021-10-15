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
}
