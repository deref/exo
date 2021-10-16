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
}
