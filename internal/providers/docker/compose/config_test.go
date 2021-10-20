package compose

import "testing"

func TestConfigYAML(t *testing.T) {
	testYAML(t, "config", `
file: ./path/to/file
external: true
name: something
`, Config{
		File:     MakeString("./path/to/file"),
		External: MakeBool(true),
		Name:     MakeString("something"),
	})
	assertInterpolated(t, map[string]string{"x": "1", "ok": "true"}, `
file: ${x}
external: ${ok}
`, Config{
		File: MakeString("${x}").WithValue("1"),
		External: Bool{
			String: MakeString("${ok}").WithValue("true"),
			Value:  true,
		},
	})
}
