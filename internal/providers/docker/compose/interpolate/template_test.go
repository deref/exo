package interpolate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	okTests := []struct {
		Template string
		Expected string
		Env      map[string]string
	}{
		{"", "", map[string]string{}},
		{"$~", "$~", map[string]string{}},            // Ignore invalid.
		{"$$", "$", map[string]string{}},             // Escaped dollar sign.
		{"$$x", "$x", map[string]string{"x": "y"}},   // Skip after escaped.
		{"$x", "", map[string]string{}},              // No substitute.
		{"${x}", "", map[string]string{}},            // No substitute, with braces.
		{"$x", "y", map[string]string{"x": "y"}},     // Simple substitute.
		{"${x}", "y", map[string]string{"x": "y"}},   // Simple substitute, with braces.
		{"${x:-y}", "y", map[string]string{}},        // Default unset.
		{"${x:-y}", "y", map[string]string{"x": ""}}, // Default empty.
		{"${x:-y}", "z", map[string]string{"x": "z"}},
		{"${x-y}", "y", map[string]string{}},       // Default unset.
		{"${x-y}", "", map[string]string{"x": ""}}, // Substitute empty.
		{"${x-y}", "z", map[string]string{"x": "z"}},

		{"${x:?some error}", "y", map[string]string{"x": "y"}},
		{"${x?some error}", "", map[string]string{"x": ""}},
		{"${x?some error}", "y", map[string]string{"x": "y"}},

		{"abc$v1 ${v2}xyz", "abcVONE VTWOxyz", map[string]string{"v1": "VONE", "v2": "VTWO"}},

		{"$X", "UPPER", map[string]string{"X": "UPPER", "x": "lower"}}, // Case-sensitive.

		{"prefix${ten}", "prefix10", map[string]string{"ten": "10"}},
		{"${ten}suffix", "10suffix", map[string]string{"ten": "10"}},
		{"prefix${ten}suffix", "prefix10suffix", map[string]string{"ten": "10"}},
	}
	for _, test := range okTests {
		tmpl, err := NewTemplate(test.Template)
		assert.NoError(t, err)
		actual, err := Substitute(tmpl, MapEnvironment(test.Env))
		assert.NoError(t, err)
		assert.Equal(t, test.Expected, actual, "test=%#v", test)
	}

	errTests := []struct {
		Template string
		Message  string
		Env      map[string]string
	}{
		{"${x:?some error}", "some error", map[string]string{}},
		{"${x:?some error}", "some error", map[string]string{"x": ""}},
	}
	for _, test := range errTests {
		tmpl, err := NewTemplate(test.Template)
		assert.NoError(t, err)
		_, err = Substitute(tmpl, MapEnvironment(test.Env))
		assert.EqualError(t, err, test.Message, "test=%#v", test)
	}
}
