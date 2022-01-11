package resolvers

import "context"

type WorkspaceResolver struct {
	Q *QueryResolver
	WorkspaceRow
}

type WorkspaceRow struct {
	ID        string  `db:"id"`
	ProjectID *string `db:"project_id"`
}

func (r *QueryResolver) WorkspaceByID(ctx context.Context, args struct {
	ID string
}) (*WorkspaceResolver, error) {
	return r.workspaceByID(ctx, &args.ID)
}

func (r *QueryResolver) workspaceByID(ctx context.Context, id *string) (*WorkspaceResolver, error) {
	ws := &WorkspaceResolver{}
	err := r.getRowByID(ctx, &ws.WorkspaceRow, `
		SELECT id, project_id
		FROM workspace
		WHERE id = ?
	`, id)
	if ws.ID == "" {
		ws = nil
	}
	return ws, err
}

func (r *WorkspaceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r.Q.projectByID(ctx, r.ProjectID)
}
