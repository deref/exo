package compose

import "testing"

func TestStringsYAML(t *testing.T) {
	testYAML(t, "string", `str`, Strings{
		Values: []string{
			"str",
		},
	})
	testYAML(t, "empty", `[]`, Strings{
		IsSequence: true,
		Values:     []string{},
	})
	testYAML(t, "single", `
- elem
`, Strings{
		IsSequence: true,
		Values: []string{
			"elem",
		},
	})
	testYAML(t, "multiple", `
- one
- two
`, Strings{
		IsSequence: true,
		Values: []string{
			"one",
			"two",
		},
	})
}
