// See also doc/reconciliation.md

package resolvers

import (
	"context"
	"fmt"
)

type ReconciliationResolver struct {
	Q         *RootResolver
	StackID   string
	Component *ComponentResolver
	Job       *JobResolver
}

// Starts a reconciliation job (top-level async task) to reconcile all components in a stack.
func (r *MutationResolver) startStackReconciliation(ctx context.Context, stack *StackResolver) (*ReconciliationResolver, error) {
	job, err := r.createJob(ctx, "reconcileStack", map[string]any{
		"ref": stack.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating reconciliation job: %w", err)
	}
	return &ReconciliationResolver{
		StackID: stack.ID,
		Job:     job,
	}, nil
}

// Starts a reconciliation job (top-level async task) to reconcile a particular component subtree.
func (r *MutationResolver) startComponentReconciliation(ctx context.Context, component *ComponentResolver) (*ReconciliationResolver, error) {
	job, err := r.createJob(ctx, "reconcileComponent", map[string]any{
		"ref": component.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating reconciliation job: %w", err)
	}
	return &ReconciliationResolver{
		StackID:   component.StackID,
		Component: component,
		Job:       job,
	}, nil
}

// Start a task within an existing reconciliation job as a step of the recursive process.
func (r *MutationResolver) startChildReconciliation(ctx context.Context, component *ComponentResolver) (*TaskResolver, error) {
	return r.createTask(ctx, "reconcileComponent", map[string]any{
		"stack": component.StackID,
		"ref":   component.ID,
	})
}

func (r *ReconciliationResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, &r.StackID)
}

func (r *ReconciliationResolver) JobID() string {
	return r.Job.ID
}

func (r *MutationResolver) ReconcileStack_label(ctx context.Context, args struct {
	Ref string
}) (string, error) {
	stack, err := r.stackByRef(ctx, &args.Ref)
	if err := validateResolve("stack", args.Ref, stack, err); err != nil {
		return "", err
	}
	return fmt.Sprintf("reconcile %s", stack.Name), nil
}

func (r *MutationResolver) ReconcileStack(ctx context.Context, args struct {
	Ref string
}) (*VoidResolver, error) {
	stack, err := r.stackByRef(ctx, &args.Ref)
	if err := validateResolve("stack", args.Ref, stack, err); err != nil {
		return nil, err
	}
	componentSet := &componentSetResolver{
		Q:       r,
		StackID: stack.ID,
		All:     true,
	}
	components, err := componentSet.Items(ctx)
	taskInputs := make([]TaskInput, len(components))
	for i, component := range components {
		taskInputs[i] = TaskInput{
			Mutation: "reconcileComponent",
			Arguments: map[string]any{
				"stack": stack.ID,
				"ref":   component.ID,
			},
		}
	}
	_, err = r.createTasks(ctx, taskInputs)
	return nil, err
}

func (r *MutationResolver) ReconcileComponent_label(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (string, error) {
	component, _ := r.componentByRef(ctx, args.Ref, args.Stack)
	if component == nil {
		return "reconciling unknown component", nil
	}
	return fmt.Sprintf("reconcile %s", component.Name), nil
}

func (r *MutationResolver) ReconcileComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (*VoidResolver, error) {
	component, err := r.componentByRef(ctx, args.Ref, args.Stack)
	if err := validateResolve("component", args.Ref, component, err); err != nil {
		return nil, err
	}
	return nil, r.reconcileComponent(ctx, component)
}

func (r *MutationResolver) reconcileComponent(ctx context.Context, component *ComponentResolver) error {
	if component.Disposed == nil {
		// XXX before children reconciled hook.
		// XXX after children reconciled hook.
	} else {
		var err error
		if component, err = r.shutdownComponent(ctx, component.ID); err != nil {
			return fmt.Errorf("shutting down: %w", err)
		}
		// TODO: On success, delete the component.
		return nil
	}

	return nil
}

func (r *MutationResolver) reconcileChildren(ctx context.Context, parent *ComponentResolver) error {
	newChildren, err := parent.render(ctx)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	oldChildren, err := parent.Children(ctx)
	if err != nil {
		return fmt.Errorf("resolving children: %w", err)
	}

	type child struct {
		ID  string
		Old *ComponentDefinition
		New *ComponentDefinition
	}
	children := make(map[string]child) // Keyed by ident.

	for _, oldChild := range oldChildren {
		def := ComponentDefinition{
			Type: oldChild.Type,
			Name: oldChild.Name,
			Key:  oldChild.Key,
			Spec: oldChild.Spec,
		}
		ident := def.Ident()
		if _, exists := children[ident]; exists {
			panic(fmt.Errorf("unexpected old component ident conflict: %q", ident))
		}
		children[ident] = child{
			ID:  oldChild.ID,
			Old: &def,
		}
	}

	newNames := make(map[string]bool)
	for _, newChild := range newChildren {
		if newNames[newChild.Name] {
			return fmt.Errorf("child name conflict: %q", newChild.Name)
		}
		newNames[newChild.Name] = true

		def := newChild
		ident := def.Ident()
		child := children[ident]
		child.New = &def
		children[ident] = child
	}

	for _, child := range children {
		var err error
		var component *ComponentResolver
		switch {
		case child.Old == nil:
			def := *child.New
			component, err = r.createComponent(ctx, parent.StackID, &parent.ID, def)
			if err != nil {
				return fmt.Errorf("creating %q: %w", def.Name, err)
			}
		case child.New == nil:
			def := *child.Old
			component, err = r.disposeComponent(ctx, child.ID)
			if err != nil {
				return fmt.Errorf("disposing %q: %w", def.Name, err)
			}
		case child.New.Spec.String() != child.Old.Spec.String():
			def := *child.New
			component, err = r.updateComponent(ctx, child.ID, def.Name, def.Spec)
			if err != nil {
				return fmt.Errorf("updating %q: %w", def.Name, err)
			}
		default:
			// Unchanged, skip reconciliation.
			continue
		}

		task, err := r.startChildReconciliation(ctx, component)
		if err != nil {
			return fmt.Errorf("starting child reconciliation task for %q: %w", component.Name, err)
		}
		// TODO: If we're doing to reconcile in a loop, cannot rely on the natural
		// structured-concurrency behavior to await child tasks; instead, need an
		// explicit await.
		_ = task
	}
	return nil
}
