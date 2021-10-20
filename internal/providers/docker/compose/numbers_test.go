package compose

import "testing"

func TestNumbersYAML(t *testing.T) {
	testYAML(t, "int", `123`, MakeInt(123))
	assertInterpolated(t, map[string]string{"one": "1"}, `${one}`, Int{
		String: MakeString("${one}").WithValue("1"),
		Value:  1,
	})
}
