package resolvers

import (
	"context"
	"errors"
)

type StackResolver struct {
	Q *QueryResolver
	StackRow
}

type StackRow struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	ProjectID   *string `db:"project_id"`
	WorkspaceID *string `db:"workspace_id"`
}

func (r *QueryResolver) StackByID(ctx context.Context, args struct {
	ID string
}) (*StackResolver, error) {
	return r.stackByID(ctx, &args.ID)
}

func (r *QueryResolver) stackByID(ctx context.Context, id *string) (*StackResolver, error) {
	s := &StackResolver{}
	err := r.getRowByID(ctx, &s.StackRow, `
		SELECT id, name, project_id, workspace_id
		FROM stack
		WHERE id = ?
	`, id)
	if s.ID == "" {
		s = nil
	}
	return s, err
}

func (r *QueryResolver) stacksByWorkspace(ctx context.Context, workspaceID string) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT stack.id, stack.name, stack.project_id, workspace_id
		FROM stack
		INNER JOIN workspace ON stack.workspace_id = workspace.id
		WHERE workspace.id = ?
		ORDER BY stack.id ASC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*StackResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &StackResolver{
			Q:        r,
			StackRow: row,
		}
	}
	return resolvers, nil
}

func (r *QueryResolver) StackByRef(ctx context.Context, args struct {
	Ref string
}) (*StackResolver, error) {
	ws, err := r.workspaceByRef(ctx, args.Ref)
	if ws == nil || err != nil {
		return nil, err
	}
	return ws.Stack(ctx)
}

func (r *StackResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED: Stack.Project resolver")
}

func (r *StackResolver) Workspace(ctx context.Context) (*WorkspaceResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED: Stack.Workspace resolver")
}

func (r *StackResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByStack(ctx, r.ID)
}
