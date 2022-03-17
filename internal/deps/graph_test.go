package deps_test

import (
	"testing"

	"github.com/deref/exo/internal/deps"
	"github.com/stretchr/testify/assert"
)

func TestImmediateDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn(deps.StringNode("x"), deps.StringNode("y")))

	assert.True(t, g.DependsOn(deps.StringNode("x"), deps.StringNode("y")))
	assert.True(t, g.HasDependent(deps.StringNode("y"), deps.StringNode("x")))
	assert.False(t, g.DependsOn(deps.StringNode("y"), deps.StringNode("x")))
	assert.False(t, g.HasDependent(deps.StringNode("x"), deps.StringNode("y")))

	// No self-dependencies.
	assert.Error(t, g.DependOn(deps.StringNode("z"), deps.StringNode("z")))
	// No bidirectional dependencies.
	assert.Error(t, g.DependOn(deps.StringNode("y"), deps.StringNode("x")))
}

func TestTransitiveDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn(deps.StringNode("x"), deps.StringNode("y")))
	assert.NoError(t, g.DependOn(deps.StringNode("y"), deps.StringNode("z")))

	assert.True(t, g.DependsOn(deps.StringNode("x"), deps.StringNode("z")))
	assert.True(t, g.HasDependent(deps.StringNode("z"), deps.StringNode("x")))
	assert.False(t, g.DependsOn(deps.StringNode("z"), deps.StringNode("x")))
	assert.False(t, g.HasDependent(deps.StringNode("x"), deps.StringNode("")))

	// No circular dependencies.
	assert.Error(t, g.DependOn(deps.StringNode("z"), deps.StringNode("x")))
}

func TestLeaves(t *testing.T) {
	g := deps.New()
	g.DependOn(deps.StringNode("cake"), deps.StringNode("eggs"))
	g.DependOn(deps.StringNode("cake"), deps.StringNode("flour"))
	g.DependOn(deps.StringNode("eggs"), deps.StringNode("chickens"))
	g.DependOn(deps.StringNode("flour"), deps.StringNode("grain"))
	g.DependOn(deps.StringNode("chickens"), deps.StringNode("feed"))
	g.DependOn(deps.StringNode("chickens"), deps.StringNode("grain"))
	g.DependOn(deps.StringNode("grain"), deps.StringNode("soil"))

	leaves := g.Leaves()
	assert.ElementsMatch(t, leaves, []deps.StringNode{"feed", "soil"})
}

func TestTopologicalSort(t *testing.T) {
	g := deps.New()
	g.DependOn(deps.StringNode("cake"), deps.StringNode("eggs"))
	g.DependOn(deps.StringNode("cake"), deps.StringNode("flour"))
	g.DependOn(deps.StringNode("eggs"), deps.StringNode("chickens"))
	g.DependOn(deps.StringNode("flour"), deps.StringNode("grain"))
	g.DependOn(deps.StringNode("chickens"), deps.StringNode("grain"))
	g.DependOn(deps.StringNode("grain"), deps.StringNode("soil"))

	sorted := g.TopoSorted()
	pairs := []struct {
		before string
		after  string
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
	comesBefore := func(before, after any) bool {
		iBefore := -1
		iAfter := -1
		for i, elem := range sorted {
			if elem.ID() == before {
				iBefore = i
			} else if elem.ID() == after {
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

	g.DependOn(deps.StringNode("web"), deps.StringNode("database"))
	g.DependOn(deps.StringNode("web"), deps.StringNode("aggregator"))
	g.DependOn(deps.StringNode("aggregator"), deps.StringNode("database"))
	g.DependOn(deps.StringNode("web"), deps.StringNode("logger"))
	g.DependOn(deps.StringNode("web"), deps.StringNode("config"))
	g.DependOn(deps.StringNode("web"), deps.StringNode("metrics"))
	g.DependOn(deps.StringNode("database"), deps.StringNode("config"))
	g.DependOn(deps.StringNode("metrics"), deps.StringNode("config"))
	/*
		   /--------------\
		web - aggregator - database
		   \_ logger               \
		    \___________________ config
		     \________ metrics _/
	*/

	layers := g.TopoSortedLayers()
	assert.Len(t, layers, 4)
	assert.ElementsMatch(t, []deps.StringNode{"config", "logger"}, layers[0])
	assert.ElementsMatch(t, []deps.StringNode{"database", "metrics"}, layers[1])
	assert.ElementsMatch(t, []deps.StringNode{"aggregator"}, layers[2])
	assert.ElementsMatch(t, []deps.StringNode{"web"}, layers[3])
}
