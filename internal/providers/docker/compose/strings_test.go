package compose

import "testing"

func TestStringsYAML(t *testing.T) {
	testYAML(t, "string", `str`, Strings{
		Items: []string{
			"str",
		},
	})
	testYAML(t, "empty", `[]`, Strings{
		IsSequence: true,
		Items:      []string{},
	})
	testYAML(t, "single", `
- elem
`, Strings{
		IsSequence: true,
		Items: []string{
			"elem",
		},
	})
	testYAML(t, "multiple", `
- one
- two
`, Strings{
		IsSequence: true,
		Items: []string{
			"one",
			"two",
		},
	})
}
