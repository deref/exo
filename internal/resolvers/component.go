package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/manifest/exocue"
	. "github.com/deref/exo/internal/scalars"
)

type ComponentResolver struct {
	Q *QueryResolver
	ComponentRow
}

type ComponentRow struct {
	ID       string     `db:"id"`
	StackID  string     `db:"stack_id"`
	ParentID *string    `db:"parent_id"`
	Type     string     `db:"type"`
	Name     string     `db:"name"`
	Key      string     `db:"key"`
	Spec     string     `db:"spec"`
	State    JSONObject `db:"state"`
	Disposed *Instant   `db:"disposed"`
}

func (r *QueryResolver) ComponentByID(ctx context.Context, args struct {
	ID string
}) (*ComponentResolver, error) {
	return r.componentByID(ctx, &args.ID)
}

func (r *QueryResolver) componentByID(ctx context.Context, id *string) (*ComponentResolver, error) {
	component := &ComponentResolver{
		Q: r,
	}
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

	component := &ComponentResolver{
		Q: r,
	}
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

func (r *QueryResolver) componentsByStack(ctx context.Context, stackID string, all bool) ([]*ComponentResolver, error) {
	var rows []ComponentRow
	var q string
	if all {
		q = `
			SELECT *
			FROM component
			WHERE stack_id = ?
			ORDER BY name ASC
		`
	} else {
		q = `
			SELECT *
			FROM component
			WHERE stack_id = ?
			AND disposed IS NULL
			ORDER BY name ASC
		`
	}
	err := r.DB.SelectContext(ctx, &rows, q, stackID)
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

func (r *QueryResolver) componentsByParent(ctx context.Context, stackID string) ([]*ComponentResolver, error) {
	var rows []ComponentRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM component
		WHERE parent_id = ?
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

func (r *ComponentResolver) Parent(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, r.ParentID)
}

func (r *ComponentResolver) Children(ctx context.Context) ([]*ComponentResolver, error) {
	return r.Q.componentsByParent(ctx, r.ID)
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

	// TODO: Validate name & spec.

	row := ComponentRow{
		ID:      gensym.RandomBase32(),
		StackID: stack.ID,
		Name:    args.Name,
		Type:    args.Type,
		Spec:    args.Spec,
		State:   make(JSONObject),
	}
	if err := r.insertRow(ctx, "component", row); err != nil {
		if isSqlConflict(err) {
			return nil, conflictErrorf("a component named %q already exists", row.Name)
		}
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return r.beginComponentReconciliation(ctx, row)
}

func (r *MutationResolver) UpdateComponent(ctx context.Context, args struct {
	Stack   *string
	Ref     string
	NewSpec *string
	NewName *string
}) (*ReconciliationResolver, error) {
	component, err := r.componentByRef(ctx, args.Ref, args.Stack)
	if err != nil {
		return nil, fmt.Errorf("resolving component: %w", err)
	}
	if component == nil {
		return nil, errors.New("no such component")
	}

	spec := component.Spec
	if args.NewSpec != nil {
		spec = *args.NewSpec
	}

	name := component.Name
	if args.NewName != nil {
		name = *args.NewName
	}

	// TODO: Validate name and spec.

	var row ComponentRow
	if _, err := r.DB.ExecContext(ctx, `
		UPDATE component
		SET spec = ?, name = ?
		WHERE id = ?
	`, spec, name, component.ID); err != nil {
		return nil, err
	}
	return r.beginComponentReconciliation(ctx, row)
}

func (r *MutationResolver) DisposeComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (*ReconciliationResolver, error) {
	component, err := r.componentByRef(ctx, args.Ref, args.Stack)
	if err != nil {
		return nil, fmt.Errorf("resolving component: %w", err)
	}
	if component == nil {
		return nil, errors.New("no such component")
	}
	now := Now(ctx)
	var row ComponentRow
	if err := r.DB.GetContext(ctx, &row, `
		UPDATE component
		SET disposed = COALESCE(disposed, ?)
		WHERE id = ?
		RETURNING *
	`, now, component.ID); err != nil {
		return nil, err
	}
	return r.beginComponentReconciliation(ctx, row)
}

func (r *MutationResolver) beginComponentReconciliation(ctx context.Context, row ComponentRow) (*ReconciliationResolver, error) {
	job, err := r.createJob(ctx, "reconcileComponent", map[string]interface{}{
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

func (r *ComponentResolver) Configuration(ctx context.Context) (string, error) {
	cfg, err := r.configuration(ctx)
	if err != nil {
		return "", err
	}
	return exocue.FormatString(exocue.StructToFile(cfg))
}

func (r *ComponentResolver) configuration(ctx context.Context) (cue.Value, error) {
	stack, err := r.Stack(ctx)
	if err != nil {
		return cue.Value{}, fmt.Errorf("resolving stack: %w", err)
	}
	cfg, err := stack.configuration(ctx)
	if err != nil {
		return cue.Value{}, err
	}
	return cfg.Component(r.Name), nil
}
