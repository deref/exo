package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/manifest/exocue"
	"github.com/deref/exo/internal/providers/sdk"
	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
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
	err = r.db.GetContext(ctx, &component.ComponentRow, `
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

func (r *QueryResolver) componentByResourceID(ctx context.Context, resourceID *string) (*ComponentResolver, error) {
	component := &ComponentResolver{
		Q: r,
	}
	err := r.getRowByKey(ctx, &component.ComponentRow, `
		SELECT *
		FROM component
		WHERE resource_id = ?
	`, resourceID)
	if component.ID == "" {
		component = nil
	}
	return component, err
}

type componentSetResolver struct {
	Q         *RootResolver
	StackID   string
	All       bool
	Recursive bool
}

func (r *componentSetResolver) items(ctx context.Context) ([]*ComponentResolver, error) {
	var rows []ComponentRow
	var q string
	// Utilizes the `component_path` index.
	q = `
		SELECT *
		FROM component
		WHERE stack_id = ?
		AND IIF(?, true, COALESCE(parent_id, stack_id) = stack_id)
		AND IIF(?, true, disposed IS NULL)
		ORDER BY parent_id, name ASC
	`
	err := r.Q.db.SelectContext(ctx, &rows, q, r.StackID, r.Recursive, r.All)
	if err != nil {
		return nil, err
	}
	return componentRowsToResolvers(r.Q, rows), nil
}

func (r *QueryResolver) componentsByStack(ctx context.Context, stackID string) ([]*ComponentResolver, error) {
	componentSet := &componentSetResolver{
		Q:       r,
		StackID: stackID,
	}
	return componentSet.items(ctx)
}

func (r *QueryResolver) componentsByParent(ctx context.Context, parentID string) ([]*ComponentResolver, error) {
	var rows []ComponentRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM component
		WHERE parent_id = ?
		AND disposed IS NULL
		ORDER BY name ASC
	`, parentID)
	if err != nil {
		return nil, err
	}
	return componentRowsToResolvers(r, rows), nil
}

func componentRowsToResolvers(r *RootResolver, rows []ComponentRow) []*ComponentResolver {
	resolvers := make([]*ComponentResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ComponentResolver{
			Q:            r,
			ComponentRow: row,
		}
	}
	return resolvers
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
		return nil, errutil.HTTPErrorf(http.StatusNotFound, "no such stack: %q", args.Stack)
	}

	row, err := r.createComponent(ctx, stack.ID /* parentID: */, nil, args.Type, args.Name, args.Spec, "")
	if err != nil {
		return nil, err
	}
	reconciliation, err := r.startComponentReconciliationJob(ctx, row)
	if err != nil {
		return nil, fmt.Errorf("starting component reconciliation: %w", err)
	}
	return reconciliation, nil
}

type NewComponentInput struct {
	Type string // XXX doesn't belong here?
	Name string
	Key  string
	Spec string
}

func (r *MutationResolver) createComponent(ctx context.Context, stackID string, parentID *string, typ, name, spec, key string) (*ComponentResolver, error) {
	// TODO: Validate type, name, key & spec.

	row := ComponentRow{
		ID:       gensym.RandomBase32(),
		StackID:  stackID,
		ParentID: parentID,
		Name:     name,
		Type:     typ,
		Key:      key,
		Spec:     spec,
		State:    make(JSONObject),
	}
	if err := r.insertRow(ctx, "component", row); err != nil {
		if isSqlConflict(err) {
			return nil, conflictErrorf("a component named %q already exists", row.Name)
		}
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &ComponentResolver{
		Q:            r,
		ComponentRow: row,
	}, nil
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
	if err := r.db.GetContext(ctx, &row, `
		UPDATE component
		SET spec = ?, name = ?
		WHERE id = ?
		RETURNING *
	`, spec, name, component.ID); err != nil {
		return nil, err
	}
	component = &ComponentResolver{
		Q:            r,
		ComponentRow: row,
	}
	return r.startComponentReconciliationJob(ctx, component)
}

func (r *MutationResolver) DisposeComponents(ctx context.Context, args struct {
	Stack *string
	Refs  []string
}) (*ReconciliationResolver, error) {
	if len(args.Refs) != 1 {
		panic("TODO: Bulk DisposeComponents")
	}
	return r.DisposeComponent(ctx, struct {
		Stack *string
		Ref   string
	}{
		Ref:   args.Refs[0],
		Stack: args.Stack,
	})
}

func (r *MutationResolver) DisposeComponent(ctx context.Context, args struct {
	Stack *string
	Ref   string
}) (*ReconciliationResolver, error) {
	// TODO: Implement in terms of DisposeComponents.
	component, err := r.componentByRef(ctx, args.Ref, args.Stack)
	if err != nil {
		return nil, fmt.Errorf("resolving component: %w", err)
	}
	if component == nil {
		return nil, errors.New("no such component")
	}
	component, err = r.disposeComponent(ctx, component.ID)
	if err != nil {
		return nil, err
	}
	return r.startComponentReconciliationJob(ctx, component)
}

func (r *MutationResolver) disposeComponent(ctx context.Context, id string) (*ComponentResolver, error) {
	now := Now(ctx)
	var row ComponentRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE component
		SET disposed = COALESCE(disposed, ?)
		WHERE id IN (
			WITH RECURSIVE rec (id) AS (
				SELECT ?
				UNION
				SELECT component.id FROM component, rec WHERE component.parent_id = rec.id
			)
			SELECT id FROM rec
		)
		RETURNING *
	`, now, id); err != nil {
		return nil, err
	}
	return &ComponentResolver{
		Q:            r,
		ComponentRow: row,
	}, nil
}

func (r *MutationResolver) startComponentReconciliationJob(ctx context.Context, component *ComponentResolver) (*ReconciliationResolver, error) {
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

func (r *ComponentResolver) Configuration(ctx context.Context, args struct {
	Recursive *bool
	Final     *bool
}) (string, error) {
	cfg, err := r.configuration(ctx, isTrue(args.Recursive))
	if err != nil {
		return "", err
	}
	return formatConfiguration(cue.Value(cfg), isTrue(args.Final))
}

func (r *ComponentResolver) configuration(ctx context.Context, recursive bool) (exocue.Component, error) {
	stack, err := r.Stack(ctx)
	if err != nil {
		return exocue.Component{}, fmt.Errorf("resolving stack: %w", err)
	}
	b := exocue.NewBuilder()
	if err := stack.addConfiguration(ctx, b, recursive); err != nil {
		return exocue.Component{}, fmt.Errorf("adding stack configuration: %w", err)
	}
	return b.Build().Component(r.ID), nil
}

func (r *ComponentResolver) controller(ctx context.Context) (*sdk.Controller, error) {
	controller := getController(ctx, r.Type)
	if controller == nil {
		return nil, fmt.Errorf("no controller for type: %q", r.Type)
	}
	return controller, nil
}

func (r *MutationResolver) ensureComponentInitialize(ctx context.Context, componentID string) error {
	return r.ensureTask(ctx, "initializeComponent", map[string]any{
		"id": componentID,
	}, componentID)
}

func (r *MutationResolver) controlComponent(ctx context.Context, id string, f func(*sdk.Controller, exocue.Component) error) error {
	component, err := r.componentByID(ctx, &id)
	if err != nil {
		return fmt.Errorf("resolving component: %w", err)
	}
	if component == nil {
		return fmt.Errorf("no such component: %q", id)
	}

	controller, err := component.controller(ctx)
	if err != nil {
		return fmt.Errorf("resolving controller: %w", err)
	}

	configuration, err := component.configuration(ctx, true)
	if err != nil {
		return fmt.Errorf("resolving configuration: %w", err)
	}

	return f(controller, configuration)
}

func (r *MutationResolver) InitializeComponent(ctx context.Context, args struct {
	ID string
}) (*VoidResolver, error) {
	err := r.controlComponent(ctx, args.ID, func(controller *sdk.Controller, configuration exocue.Component) error {
		return controller.Initialize(ctx, cue.Value(configuration))
	})
	return nil, err
}

func (r *MutationResolver) ensureComponentTransition(ctx context.Context, id, spec string) error {
	return r.ensureTask(ctx, "transitionComponent", map[string]any{
		"id":   id,
		"spec": spec,
	}, id)
}

func (r *MutationResolver) TransitionComponent(ctx context.Context, args struct {
	ID   string
	Spec string
}) (*VoidResolver, error) {
	return nil, errors.New("TODO: transition component")
}

func (r *MutationResolver) ensureComponentShutdown(ctx context.Context, componentID string) error {
	return r.ensureTask(ctx, "shutdownComponent", map[string]any{
		"id": componentID,
	}, componentID)
}

func (r *MutationResolver) ShutdownComponent(ctx context.Context, args struct {
	ID string
}) (*VoidResolver, error) {
	// XXX if there are still children, abort and try again later.
	// after done, trigger reconciliation of parent.
	// ^^^ actually, this doesn't make sense, the parent reconcilliation should wait?
	return nil, errors.New("TODO: shutdown component")
}

func (r *ComponentResolver) Reconciling() bool {
	return false // XXX
}

func (r *ComponentResolver) Running() bool {
	return true // XXX
}

func (r *ComponentResolver) AsProcess(ctx context.Context) *ProcessComponentResolver {
	return r.Q.processFromComponent(r)
}

func (r *ComponentResolver) AsStore(ctx context.Context) *StoreComponentResolver {
	return r.Q.storeFromComponent(r)
}
