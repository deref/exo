package compose

import "testing"

func TestBoolYAML(t *testing.T) {
	testYAML(t, "true", `true`, MakeBool(true))
	testYAML(t, "false", `false`, MakeBool(false))
	assertInterpolated(t, map[string]string{"yup": "true"}, `${yup}`, Bool{
		String: MakeString("${yup}").WithValue("true"),
		Value:  true,
	})
	assertInterpolated(t, map[string]string{"nope": "false"}, `${nope}`, Bool{
		String: MakeString("${nope}").WithValue("false"),
		Value:  false,
	})
}
