package compose

import "testing"

func TestNumbersYAML(t *testing.T) {
	testYAML(t, "int", `123`, MakeInt(123))
}
