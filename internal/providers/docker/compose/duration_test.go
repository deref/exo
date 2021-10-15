package compose

import (
	"testing"
	"time"
)

func TestDurationYAML(t *testing.T) {
	check := func(s string, expected time.Duration) {
		testYAML(t, s, s, Duration(expected))
	}
	check("5s", 5*time.Second)
	check("1m30s", 90*time.Second)
	check("37us", 37*time.Microsecond)
	check("1h5m30s20ms", 1*time.Hour+5*time.Minute+30*time.Second+20*time.Millisecond)
}
