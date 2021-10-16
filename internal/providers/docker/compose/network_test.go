package compose

import "testing"

func TestNetworkYAML(t *testing.T) {
	assertInterpolated(t, map[string]string{"x": "y"}, `
driver: ${x}
`, Network{
		Driver: MakeString("${x}").WithValue("y"),
	})
}
