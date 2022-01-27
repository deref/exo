package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type ComponentResolver struct {
	Q *QueryResolver
	ComponentRow
}

type ComponentRow struct {
	ID       string   `db:"id"`
	StackID  string   `db:"stack_id"`
	Name     string   `db:"name"`
	Type     string   `db:"type"`
	Spec     string   `db:"spec"`
	Disposed *Instant `db:"disposed"`
}

func (r *QueryResolver) ComponentByID(ctx context.Context, args struct {
	ID string
}) (*ComponentResolver, error) {
	return r.componentByID(ctx, &args.ID)
}

func (r *QueryResolver) componentByID(ctx context.Context, id *string) (*ComponentResolver, error) {
	component := &ComponentResolver{}
	err := r.getRowByKey(ctx, &component.ComponentRow, `
		SELECT *
		FROM component
		WHERE id = ?
	`, id)
	if component.ID == "" {
		component = nil
	}
	return component, err
}

func (r *QueryResolver) componentByName(ctx context.Context, stack string, name string) (*ComponentResolver, error) {
	stackResolver, err := r.stackByRef(ctx, &stack)
	if stackResolver == nil || err != nil {
		return nil, err
	}
	stackID := stackResolver.ID

	component := &ComponentResolver{}
	err = r.DB.GetContext(ctx, &component.ComponentRow, `
		SELECT *
		FROM component
		WHERE stack_id = ?
		AND name = ?
	`, stackID, name)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if component.ID == "" {
		component = nil
	}
	return component, err
}

func (r *QueryResolver) ComponentByRef(ctx context.Context, args struct {
	Ref   string
	Stack *string
}) (*ComponentResolver, error) {
	return r.componentByRef(ctx, args.Ref, args.Stack)
}

func (r *QueryResolver) componentByRef(ctx context.Context, ref string, stack *string) (*ComponentResolver, error) {
	component, err := r.componentByID(ctx, &ref)
	if component != nil || err != nil {
		return component, err
	}
	if stack != nil {
		component, err = r.componentByName(ctx, *stack, ref)
	}
	return component, err
}

func (r *QueryResolver) componentsByStack(ctx context.Context, stackID string) ([]*ComponentResolver, error) {
	var rows []ComponentRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM component
		WHERE stack_id = ?
		AND disposed IS NULL
		ORDER BY name ASC
	`, stackID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ComponentResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ComponentResolver{
			Q:            r,
			ComponentRow: row,
		}
	}
	return resolvers, nil
}

func (r *ComponentResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, &r.StackID)
}

func (r *ComponentResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByComponent(ctx, r.ID)
}

func (r *MutationResolver) CreateComponent(ctx context.Context, args struct {
	Stack string
	Name  string
	Type  string
	Spec  string
}) (*ReconciliationResolver, error) {
	stack, err := r.stackByRef(ctx, &args.Stack)
	if err != nil {
		return nil, fmt.Errorf("resolving stack: %w", err)
	}
	if stack == nil {
		return nil, fmt.Errorf("no such stack: %q", args.Stack)
	}

	row := ComponentRow{
		ID:      gensym.RandomBase32(),
		StackID: stack.ID,
		Name:    args.Name,
		Type:    args.Type,
		Spec:    args.Spec,
	}
	if err := r.insertRow(ctx, "component", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return r.beginComponentReconciliation(ctx, row)
}

func (r *MutationResolver) UpdateComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
	Spec  string
}) (*ReconciliationResolver, error) {
	return nil, fmt.Errorf("TODO: UpdateComponent")
}

func (r *MutationResolver) DisposeComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (*ReconciliationResolver, error) {
	return nil, fmt.Errorf("TODO: DisposeComponent")
}

func (r *MutationResolver) beginComponentReconciliation(ctx context.Context, row ComponentRow) (*ReconciliationResolver, error) {
	job, err := r.createJob(ctx, newTaskID(), "reconcileComponent", map[string]interface{}{
		"ref": row.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating reconciliation job: %w", err)
	}
	return &ReconciliationResolver{
		Component: &ComponentResolver{
			Q:            r,
			ComponentRow: row,
		},
		Job: job,
	}, nil
}
