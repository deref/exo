package compose

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStringYAML(t *testing.T) {
	testYAML(t, "single_quoted", `'abc'`, String{
		Style:      yaml.SingleQuotedStyle,
		Expression: "abc",
		Value:      "abc",
	})
	testYAML(t, "double_quoted", `"abc"`, String{
		Style:      yaml.DoubleQuotedStyle,
		Expression: "abc",
		Value:      "abc",
	})
	testYAML(t, "plain", `abc`, String{
		Expression: "abc",
		Value:      "abc",
	})
}
