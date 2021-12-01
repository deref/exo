package server

import (
	"github.com/deref/exo/internal/core/api"
)

type dependencyOrder int

const (
	dependencyOrderNatural dependencyOrder = iota
	dependencyOrderReverse
)

type componentQuery struct {
	DependencyOrder     dependencyOrder
	Refs                []string
	Types               []string
	IncludeDependencies bool
	IncludeDependents   bool
}

type componentQueryUpdate func(*componentQuery)

func makeComponentQuery(updates ...componentQueryUpdate) componentQuery {
	q := &componentQuery{}

	for _, update := range updates {
		update(q)
	}

	return *q
}

func withTypes(types ...string) componentQueryUpdate {
	return componentQueryUpdate(func(q *componentQuery) {
		q.Types = types
	})
}

func withRefs(refs ...string) componentQueryUpdate {
	return componentQueryUpdate(func(q *componentQuery) {
		q.Refs = refs
	})
}

var withReversedDependencies = componentQueryUpdate(func(q *componentQuery) {
	q.DependencyOrder = dependencyOrderReverse
})
var withDependencies = componentQueryUpdate(func(q *componentQuery) {
	q.IncludeDependencies = true
})
var withDependents = componentQueryUpdate(func(q *componentQuery) {
	q.IncludeDependents = true
})

var allComponentsQuery = makeComponentQuery()

var processTypes = []string{"process", "container"}
var runnableTypes = append(processTypes, "apigateway")

func isRunnableType(name string) bool {
	for _, typ := range runnableTypes {
		if name == typ {
			return true
		}
	}
	return false
}

func allProcessQuery(updates ...componentQueryUpdate) componentQuery {
	updates = append([]componentQueryUpdate{withTypes(processTypes...)}, updates...)
	return makeComponentQuery(updates...)
}

func allRunnableQuery(updates ...componentQueryUpdate) componentQuery {
	updates = append([]componentQueryUpdate{withTypes(runnableTypes...)}, updates...)
	return makeComponentQuery(updates...)
}

func allBuildableQuery(updates ...componentQueryUpdate) componentQuery {
	updates = append([]componentQueryUpdate{withTypes("container")}, updates...)
	return makeComponentQuery(updates...)
}

func (q componentQuery) describeComponentsInput(ws *Workspace) *api.DescribeComponentsInput {
	return &api.DescribeComponentsInput{
		Refs:                q.Refs,
		Types:               q.Types,
		IncludeDependencies: q.IncludeDependencies,
		IncludeDependents:   q.IncludeDependents,
	}
}
