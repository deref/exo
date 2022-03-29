package resolvers

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/util/hashutil"
	"github.com/deref/exo/internal/util/jsonutil"
)

type ReconciliationResolver struct {
	Q         *RootResolver
	StackID   string
	Component *ComponentResolver
	Job       *JobResolver
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

	// TODO: do not give up on first error; process all children.

	if component.Disposed != nil {
		if err := r.ensureComponentShutdown(ctx, component.ID); err != nil {
			return nil, fmt.Errorf("ensuring %q shutdown: %w", component.Name, err)
		}
	} else {
		controller, err := component.controller(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving controller: %w", err)
		}

		configuration, err := component.configuration(ctx, true)
		if err != nil {
			return nil, fmt.Errorf("resolving configuration: %w", err)
		}

		oldChildren, err := component.Children(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving children: %w", err)
		}

		newChildren, err := controller.Render(ctx, cue.Value(configuration))
		if err != nil {
			return nil, fmt.Errorf("rendering: %w", err)
		}

		type Child struct {
			Type     string
			ID       string
			Name     string
			Key      string
			Existing bool
			NewSpec  *string
		}
		parent := component
		names := make(map[string]bool)
		children := make(map[string]Child)
		for _, oldChild := range oldChildren {
			ident := makeChildComponentIdent(oldChild.Type, oldChild.Name, oldChild.Spec, oldChild.Key)
			child := children[ident]
			child.Type = oldChild.Type
			child.Name = oldChild.Name
			child.Key = oldChild.Key
			child.Existing = true
			children[ident] = child
		}
		for _, newChild := range newChildren {
			if names[newChild.Name] {
				return nil, fmt.Errorf("child name conflict: %q", newChild.Name)
			}
			names[newChild.Name] = true

			spec, err := jsonutil.MarshalString(newChild.Spec)
			if err != nil {
				return nil, fmt.Errorf("marshaling child %q spec: %w", newChild.Name, err)
			}

			ident := makeChildComponentIdent(newChild.Type, newChild.Name, spec, newChild.Key)
			child := children[ident]
			child.Type = newChild.Type
			child.Name = newChild.Name
			child.Key = newChild.Key
			child.NewSpec = &spec
			children[ident] = child
		}

		for _, child := range children {
			reconcileID := child.ID
			switch {
			case !child.Existing:
				spec := *child.NewSpec
				created, err := r.createComponent(ctx, parent.StackID, &parent.ID, child.Type, child.Name, spec, child.Key)
				if err != nil {
					return nil, fmt.Errorf("creating %q: %w", child.Name, err)
				}
				reconcileID = created.ID
			case child.NewSpec == nil:
				if _, err := r.disposeComponent(ctx, child.ID); err != nil {
					return nil, fmt.Errorf("disposing %q: %w", child.Name, err)
				}
			default:
				// XXX instead of task keys, might be better to have prev/next specs on
				// components, then iif changing the next step, trigger a transition. but
				// need to be concerned with concurrent changes.
				spec := *child.NewSpec
				if err := r.ensureComponentTransition(ctx, child.ID, spec); err != nil {
					return nil, fmt.Errorf("ensuring %q transition: %w", child.Name, err)
				}
			}
			if reconcileID != "" {
				_, err := r.createTask(ctx, "reconcileComponent", map[string]any{
					"ref": reconcileID,
				})
				if err != nil {
					return nil, fmt.Errorf("creating child %q reconciliation task: %w", child.Name, err)
				}
			}
		}

		if err := r.ensureComponentInitialize(ctx, component.ID); err != nil {
			return nil, fmt.Errorf("ensuring %q initialization: %w", component.Name, err)
		}
	}

	return nil, nil
}

func makeChildComponentIdent(typ, name, spec, key string) string {
	if key == "" {
		key = hashutil.Sha256Hex([]byte(spec))
	}
	return fmt.Sprintf("%s:%s:%s", typ, name, key)
}
