package compose

import (
	"testing"
)

func TestCommandYAML(t *testing.T) {
	testYAML(t, "shell", `x y z`, Command{
		IsShellForm: true,
		Parts:       []String{MakeString(`x y z`)},
	})
	testYAML(t, "parsed", `["x", "y z"]`, Command{
		Parts: []String{MakeString("x"), MakeString("y z")},
	})
}
