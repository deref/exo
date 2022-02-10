package resolvers

import (
	"context"
	"errors"
	"fmt"
)

type ReconciliationResolver struct {
	Component *ComponentResolver
	Job       *JobResolver
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

	ctrl := getController(ctx, component.Type)
	if ctrl == nil {
		return nil, fmt.Errorf("no controller for type: %q", component.Type)
	}

	cfg, err := component.configuration(ctx)

	rendered, err := ctrl.Render(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("rendering: %w", err)
	}

	fmt.Printf("!!! RENDERED: %#v\n", rendered)

	return nil, nil
}
