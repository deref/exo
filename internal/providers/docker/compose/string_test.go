package compose

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStringYAML(t *testing.T) {
	testYAML(t, "single_quoted", `'abc'`, String{
		Tag:        "!!str",
		Style:      yaml.SingleQuotedStyle,
		Expression: "abc",
		Value:      "abc",
	})
	testYAML(t, "double_quoted", `"abc"`, String{
		Tag:        "!!str",
		Style:      yaml.DoubleQuotedStyle,
		Expression: "abc",
		Value:      "abc",
	})
	testYAML(t, "plain", `abc`, String{
		Tag:        "!!str",
		Expression: "abc",
		Value:      "abc",
	})
	testYAML(t, "int", `123`, String{
		Tag:        "!!int",
		Expression: "123",
		Value:      "123",
	})
}
