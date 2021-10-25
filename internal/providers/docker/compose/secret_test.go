package compose

import "testing"

func TestSecretAML(t *testing.T) {
	assertInterpolated(t, map[string]string{"x": "y"}, `
file: ${x}
`, Secret{
		File: MakeString("${x}").WithValue("y"),
	})
}
