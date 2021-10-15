package compose

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDurationYAML(t *testing.T) {
	checkOneWay := func(expected string, d time.Duration) {
		actual, err := yaml.Marshal(Duration{
			Duration: d,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, expected, string(bytes.TrimSpace(actual)))
		}
	}
	checkRoundTrip := func(s string, d time.Duration) {
		testYAML(t, s, s, Duration{
			String:   MakeString(s),
			Duration: d,
		})
		checkOneWay(s, d)
	}
	checkRoundTrip("5s", 5*time.Second)
	checkRoundTrip("1m30s", 90*time.Second)
	checkRoundTrip("37us", 37*time.Microsecond)
	checkRoundTrip("1h5m30s20ms", 1*time.Hour+5*time.Minute+30*time.Second+20*time.Millisecond)
}
