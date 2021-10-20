package compose

import "testing"

func TestVolumeYAML(t *testing.T) {
	assertInterpolated(t, map[string]string{"x": "y"}, `
driver: ${x}
`, Volume{
		Driver: MakeString("${x}").WithValue("y"),
	})
}
