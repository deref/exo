package compose

import "testing"

func TestTupleYAML(t *testing.T) {
	testYAML(t, "string", `str`, Tuple{
		Items: []string{
			"str",
		},
	})
	testYAML(t, "empty", `[]`, Tuple{
		IsSequence: true,
		Items:      []string{},
	})
	testYAML(t, "single", `
- elem
`, Tuple{
		IsSequence: true,
		Items: []string{
			"elem",
		},
	})
	testYAML(t, "multiple", `
- one
- two
`, Tuple{
		IsSequence: true,
		Items: []string{
			"one",
			"two",
		},
	})
}
