package server

import (
	"context"
	"fmt"

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

func allProcessQuery(updates ...componentQueryUpdate) componentQuery {
	updates = append([]componentQueryUpdate{withTypes("process", "container")}, updates...)
	return makeComponentQuery(updates...)
}

func allBuildableQuery(updates ...componentQueryUpdate) componentQuery {
	updates = append([]componentQueryUpdate{withTypes("container")}, updates...)
	return makeComponentQuery(updates...)
}

func (q componentQuery) describeComponentsInput(ctx context.Context, ws *Workspace) (*api.DescribeComponentsInput, error) {
	describe := &api.DescribeComponentsInput{
		Types:               q.Types,
		IncludeDependencies: q.IncludeDependencies,
		IncludeDependents:   q.IncludeDependents,
	}

	if q.Refs != nil {
		ids, err := ws.resolveRefs(ctx, q.Refs)
		if err != nil {
			return nil, fmt.Errorf("resolving refs: %w", err)
		}
		describe.IDs = ids
	}

	return describe, nil
}
