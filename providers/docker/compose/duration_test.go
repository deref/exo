package compose

import (
	"strings"
	"testing"
	"time"

	"github.com/deref/exo/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	check := func(s string, expected time.Duration) {
		actual, err := ParseDuration(s)
		assert.NoError(t, err)
		assert.Equal(t, Duration(expected), actual)
		assert.Equal(t, s, actual.String())
	}
	check("5s", 5*time.Second)
	check("1m30s", 90*time.Second)
	check("37us", 37*time.Microsecond)
	check("1h5m30s20ms", 1*time.Hour+5*time.Minute+30*time.Second+20*time.Millisecond)
}

func TestDurationYaml(t *testing.T) {
	var actual struct {
		D Duration
	}
	yamlutil.MustUnmarshalString(`
d: 1m30s
`, &actual)
	assert.Equal(t, Duration(90*time.Second), actual.D)

	assert.Equal(t, "1h5m30s20ms", strings.TrimSpace(yamlutil.MustMarshalString(Duration(1*time.Hour+5*time.Minute+30*time.Second+20*time.Millisecond))))
}
