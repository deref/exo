package compose

import "testing"

func TestBoolYAML(t *testing.T) {
	testYAML(t, "true", `true`, MakeBool(true))
	testYAML(t, "false", `false`, MakeBool(false))
}
