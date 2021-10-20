package compose

import "testing"

func TestLoggingYAML(t *testing.T) {
	assertInterpolated(t, map[string]string{"x": "y"}, `
driver: ${x}
`, Logging{
		Driver: MakeString("${x}").WithValue("y"),
	})
}
