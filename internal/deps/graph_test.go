package deps_test

import (
	"testing"

	"github.com/deref/exo/internal/deps"
	"github.com/stretchr/testify/assert"
)

func TestImmediateDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn("x", "y"))

	assert.True(t, g.DependsOn("x", "y"))
	assert.True(t, g.HasDependent("y", "x"))
	assert.False(t, g.DependsOn("y", "x"))
	assert.False(t, g.HasDependent("x", "y"))

	// No self-dependencies.
	assert.Error(t, g.DependOn("z", "z"))
	// No bidirectional dependencies.
	assert.Error(t, g.DependOn("y", "x"))
}

func TestTransitiveDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn("x", "y"))
	assert.NoError(t, g.DependOn("y", "z"))

	assert.True(t, g.DependsOn("x", "z"))
	assert.True(t, g.HasDependent("z", "x"))
	assert.False(t, g.DependsOn("z", "x"))
	assert.False(t, g.HasDependent("x", ""))

	// No circular dependencies.
	assert.Error(t, g.DependOn("z", "x"))
}

func TestLeaves(t *testing.T) {
	g := deps.New()
	g.DependOn("cake", "eggs")
	g.DependOn("cake", "flour")
	g.DependOn("eggs", "chickens")
	g.DependOn("flour", "grain")
	g.DependOn("chickens", "feed")
	g.DependOn("chickens", "grain")
	g.DependOn("grain", "soil")

	leaves := g.Leaves()
	assert.ElementsMatch(t, leaves, []interface{}{"feed", "soil"})
}

func TestTopologicalSort(t *testing.T) {
	g := deps.New()
	g.DependOn("cake", "eggs")
	g.DependOn("cake", "flour")
	g.DependOn("eggs", "chickens")
	g.DependOn("flour", "grain")
	g.DependOn("chickens", "grain")
	g.DependOn("grain", "soil")

	sorted := g.TopoSorted()
	pairs := []struct {
		before interface{}
		after  interface{}
	}{
		{
			before: "soil",
			after:  "grain",
		},
		{
			before: "grain",
			after:  "chickens",
		},
		{
			before: "grain",
			after:  "flour",
		},
		{
			before: "chickens",
			after:  "eggs",
		},
		{
			before: "flour",
			after:  "cake",
		},
		{
			before: "eggs",
			after:  "cake",
		},
	}
	comesBefore := func(before, after interface{}) bool {
		iBefore := -1
		iAfter := -1
		for i, elem := range sorted {
			if elem == before {
				iBefore = i
			} else if elem == after {
				iAfter = i
			}
		}
		return iBefore < iAfter
	}
	for _, pair := range pairs {
		assert.True(t, comesBefore(pair.before, pair.after))
	}
}

func TestLayeredTopologicalSort(t *testing.T) {
	g := deps.New()

	g.DependOn("web", "database")
	g.DependOn("web", "aggregator")
	g.DependOn("aggregator", "database")
	g.DependOn("web", "logger")
	g.DependOn("web", "config")
	g.DependOn("web", "metrics")
	g.DependOn("database", "config")
	g.DependOn("metrics", "config")
	/*
		   /--------------\
		web - aggregator - database
		   \_ logger               \
		    \___________________ config
		     \________ metrics _/
	*/

	layers := g.TopoSortedLayers()
	assert.Len(t, layers, 4)
	assert.ElementsMatch(t, []interface{}{"config", "logger"}, layers[0])
	assert.ElementsMatch(t, []interface{}{"database", "metrics"}, layers[1])
	assert.ElementsMatch(t, []interface{}{"aggregator"}, layers[2])
	assert.ElementsMatch(t, []interface{}{"web"}, layers[3])
}
