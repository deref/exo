package telemetry

import (
	"math"
	"sort"
	"strings"
	"sync"
)

type Tags map[string]string

func NewSummaryGauge(tagLabels []string) *SummaryGauge {
	return &SummaryGauge{
		tagLabels: tagLabels,
		buckets:   make(map[string]taggedSummaryBucket),
	}
}

func (g *SummaryGauge) Observe(tags Tags, amount float64) {
	orderedTagVals := make([]string, len(g.tagLabels))
	for i, tagLabel := range g.tagLabels {
		if tagVal, ok := tags[tagLabel]; ok {
			// TODO: Escape null bytes.
			orderedTagVals[i] = tagVal
		} else {
			// Empty value sentinel.
			orderedTagVals[i] = "\xff"
		}
	}
	key := strings.Join(orderedTagVals, "\x00")

	g.mu.Lock()
	defer g.mu.Unlock()

	taggedBucket, ok := g.buckets[key]
	if !ok {
		taggedBucket = taggedSummaryBucket{
			tagLabels:     g.tagLabels,
			tagValues:     orderedTagVals,
			summaryBucket: &summaryBucket{},
		}
		g.buckets[key] = taggedBucket
	}

	taggedBucket.summaryBucket.observe(amount)
}

func (g *SummaryGauge) Buckets() []taggedSummaryBucket {
	out := make([]taggedSummaryBucket, len(g.buckets))

	var i int
	for _, bucket := range g.buckets {
		out[i] = bucket
		i++
	}

	return out
}

type SummaryGauge struct {
	tagLabels []string

	mu      sync.Mutex
	buckets map[string]taggedSummaryBucket
}

type taggedSummaryBucket struct {
	*summaryBucket
	tagLabels []string
	tagValues []string
}

func (b *taggedSummaryBucket) TagValues() []string {
	return b.tagValues
}

func (b *taggedSummaryBucket) Tags() Tags {
	tags := make(Tags, len(b.tagLabels))
	for i, tagLabel := range b.tagLabels {
		tagValue := b.tagValues[i]
		tags[tagLabel] = tagValue
	}

	return tags
}

func (b *taggedSummaryBucket) HasTagValues(tagValues []string) bool {
	if len(tagValues) != len(b.tagValues) {
		return false
	}

	for i, tv := range b.tagValues {
		if tv != tagValues[i] {
			return false
		}
	}

	return true
}

// summaryBucket keeps track of a number of observations then generates
// `SummaryStatistics` from them.
type summaryBucket struct {
	mu           sync.Mutex
	observations []float64
}

func (b *summaryBucket) observe(amount float64) {
	b.mu.Lock()
	b.observations = append(b.observations, amount)
	b.mu.Unlock()
}

func (b *summaryBucket) Summarize() SummaryStatistics {
	b.mu.Lock()
	defer b.mu.Unlock()

	stats := SummaryStatistics{
		Count: len(b.observations),
	}

	if stats.Count == 0 {
		return stats
	}

	stats.Min = math.Inf(1)
	stats.Max = math.Inf(-1)
	for _, amount := range b.observations {
		stats.Sum += amount
		stats.Min = math.Min(stats.Min, amount)
		stats.Max = math.Max(stats.Max, amount)
	}
	stats.Mean = stats.Sum / float64(stats.Count)

	// Calculate median.
	sort.Float64s(b.observations)
	middle := stats.Count / 2
	if stats.Count%2 == 0 {
		low := b.observations[middle-1]
		high := b.observations[middle]
		stats.Median = (low + high) / 2
	} else {
		stats.Median = b.observations[middle]
	}

	// Calculate standard deviation.
	var sumSqDiffs float64
	for _, amount := range b.observations {
		diff := amount - stats.Mean
		sumSqDiffs += diff * diff
	}
	variance := sumSqDiffs / float64(stats.Count)
	stats.StdDev = math.Sqrt(variance)

	return stats
}

type SummaryStatistics struct {
	Count  int     `json:"count"`
	Sum    float64 `json:"sum"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Median float64 `json:"median"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"stdDev"`
}
