package resolvers

import (
	"context"
	"errors"
	"fmt"
)

type ReconciliationResolver struct {
	Component *ComponentResolver
	Job       *TaskResolver
}

func (r *MutationResolver) ReconcileComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (*VoidResolver, error) {
	component, err := r.componentByRef(ctx, args.Ref, args.Stack)
	if err != nil {
		return nil, fmt.Errorf("resolving component: %w", err)
	}
	if component == nil {
		return nil, errors.New("no such component")
	}
	// XXX reconcile here.
	return nil, nil
}
