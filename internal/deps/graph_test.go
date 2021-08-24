package deps_test

import (
	"testing"

	"github.com/deref/exo/internal/deps"
	"github.com/stretchr/testify/assert"
)

type stringNode string

func (s stringNode) ID() string {
	return string(s)
}

func TestImmediateDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn(stringNode("x"), stringNode("y")))

	assert.True(t, g.DependsOn(stringNode("x"), stringNode("y")))
	assert.True(t, g.HasDependent(stringNode("y"), stringNode("x")))
	assert.False(t, g.DependsOn(stringNode("y"), stringNode("x")))
	assert.False(t, g.HasDependent(stringNode("x"), stringNode("y")))

	// No self-dependencies.
	assert.Error(t, g.DependOn(stringNode("z"), stringNode("z")))
	// No bidirectional dependencies.
	assert.Error(t, g.DependOn(stringNode("y"), stringNode("x")))
}

func TestTransitiveDependencies(t *testing.T) {
	g := deps.New()

	assert.NoError(t, g.DependOn(stringNode("x"), stringNode("y")))
	assert.NoError(t, g.DependOn(stringNode("y"), stringNode("z")))

	assert.True(t, g.DependsOn(stringNode("x"), stringNode("z")))
	assert.True(t, g.HasDependent(stringNode("z"), stringNode("x")))
	assert.False(t, g.DependsOn(stringNode("z"), stringNode("x")))
	assert.False(t, g.HasDependent(stringNode("x"), stringNode("")))

	// No circular dependencies.
	assert.Error(t, g.DependOn(stringNode("z"), stringNode("x")))
}

func TestLeaves(t *testing.T) {
	g := deps.New()
	g.DependOn(stringNode("cake"), stringNode("eggs"))
	g.DependOn(stringNode("cake"), stringNode("flour"))
	g.DependOn(stringNode("eggs"), stringNode("chickens"))
	g.DependOn(stringNode("flour"), stringNode("grain"))
	g.DependOn(stringNode("chickens"), stringNode("feed"))
	g.DependOn(stringNode("chickens"), stringNode("grain"))
	g.DependOn(stringNode("grain"), stringNode("soil"))

	leaves := g.Leaves()
	assert.ElementsMatch(t, leaves, []stringNode{"feed", "soil"})
}

func TestTopologicalSort(t *testing.T) {
	g := deps.New()
	g.DependOn(stringNode("cake"), stringNode("eggs"))
	g.DependOn(stringNode("cake"), stringNode("flour"))
	g.DependOn(stringNode("eggs"), stringNode("chickens"))
	g.DependOn(stringNode("flour"), stringNode("grain"))
	g.DependOn(stringNode("chickens"), stringNode("grain"))
	g.DependOn(stringNode("grain"), stringNode("soil"))

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
	comesBefore := func(before, after interface{}) bool {
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

	g.DependOn(stringNode("web"), stringNode("database"))
	g.DependOn(stringNode("web"), stringNode("aggregator"))
	g.DependOn(stringNode("aggregator"), stringNode("database"))
	g.DependOn(stringNode("web"), stringNode("logger"))
	g.DependOn(stringNode("web"), stringNode("config"))
	g.DependOn(stringNode("web"), stringNode("metrics"))
	g.DependOn(stringNode("database"), stringNode("config"))
	g.DependOn(stringNode("metrics"), stringNode("config"))
	/*
		   /--------------\
		web - aggregator - database
		   \_ logger               \
		    \___________________ config
		     \________ metrics _/
	*/

	layers := g.TopoSortedLayers()
	assert.Len(t, layers, 4)
	assert.ElementsMatch(t, []stringNode{"config", "logger"}, layers[0])
	assert.ElementsMatch(t, []stringNode{"database", "metrics"}, layers[1])
	assert.ElementsMatch(t, []stringNode{"aggregator"}, layers[2])
	assert.ElementsMatch(t, []stringNode{"web"}, layers[3])
}
