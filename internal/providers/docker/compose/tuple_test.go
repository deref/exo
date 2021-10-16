package compose

import "testing"

func TestTupleYAML(t *testing.T) {
	testYAML(t, "string", `str`, Tuple{
		Items: []String{
			MakeString("str"),
		},
	})
	testYAML(t, "empty", `[]`, Tuple{
		IsSequence: true,
		Items:      []String{},
	})
	testYAML(t, "single", `
- elem
`, Tuple{
		IsSequence: true,
		Items: []String{
			MakeString("elem"),
		},
	})
	testYAML(t, "multiple", `
- one
- two
`, Tuple{
		IsSequence: true,
		Items: []String{
			MakeString("one"),
			MakeString("two"),
		},
	})
	assertInterpolated(t, map[string]string{"x": "1"}, "${x}", Tuple{
		Items: []String{MakeString("${x}").WithValue("1")},
	})
	assertInterpolated(t, map[string]string{"x": "1"}, `
- ${x}
- y
`, Tuple{
		IsSequence: true,
		Items: []String{
			MakeString("${x}").WithValue("1"),
			MakeString("y"),
		},
	})
}
