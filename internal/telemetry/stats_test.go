package telemetry_test

import (
	"testing"

	"github.com/deref/exo/internal/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestSummaryGaugeNoTags(t *testing.T) {
	g := telemetry.NewSummaryGauge(nil)
	g.Observe(nil, -5)
	g.Observe(nil, 1)
	g.Observe(nil, 8)
	g.Observe(nil, 7)
	g.Observe(nil, 2)

	buckets := g.Buckets()
	assert.Len(t, buckets, 1)

	stats := buckets[0].Summarize()

	assert.Equal(t, 5, stats.Count)
	assert.Equal(t, float64(13), stats.Sum)
	assert.Equal(t, float64(-5), stats.Min)
	assert.Equal(t, float64(8), stats.Max)
	assert.Equal(t, float64(2), stats.Median)
	assert.InDelta(t, 2.6, stats.Mean, 0.0000001)
	assert.InDelta(t, 4.673328578, stats.StdDev, 0.000000001)
}

func TestSummaryGaugeWithTags(t *testing.T) {
	g := telemetry.NewSummaryGauge([]string{"path", "status"})
	g.Observe(telemetry.Tags{
		"path":   "/about",
		"status": "200",
	}, 100)
	g.Observe(telemetry.Tags{
		"path":   "/",
		"status": "200",
	}, 600)
	g.Observe(telemetry.Tags{
		"path":   "/about",
		"status": "200",
	}, 300)
	g.Observe(telemetry.Tags{
		"path":   "/about",
		"status": "401",
	}, 25)

	buckets := g.Buckets()
	assert.Len(t, buckets, 3)

	for _, bucket := range buckets {
		if bucket.HasTagValues([]string{"/about", "200"}) {
			assert.Equal(t, telemetry.Tags{
				"path":   "/about",
				"status": "200",
			}, bucket.Tags())
			assert.Equal(t, float64(200), bucket.Summarize().Mean)
			continue
		}

		if bucket.HasTagValues([]string{"/", "200"}) {
			assert.Equal(t, telemetry.Tags{
				"path":   "/",
				"status": "200",
			}, bucket.Tags())
			assert.Equal(t, float64(600), bucket.Summarize().Mean)
			continue
		}

		if bucket.HasTagValues([]string{"/about", "401"}) {
			assert.Equal(t, telemetry.Tags{
				"path":   "/about",
				"status": "401",
			}, bucket.Tags())
			assert.Equal(t, float64(25), bucket.Summarize().Mean)
			continue
		}
	}
}
