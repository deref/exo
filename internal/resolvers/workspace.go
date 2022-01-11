package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type WorkspaceResolver struct {
	Q *QueryResolver
	WorkspaceRow
}

type WorkspaceRow struct {
	ID        string  `db:"id"`
	Root      string  `db:"root"`
	ProjectID *string `db:"project_id"`
}

func (r *QueryResolver) AllWorkspaces(ctx context.Context) ([]*WorkspaceResolver, error) {
	var rows []WorkspaceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT id, project_id, root
		FROM workspace
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*WorkspaceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &WorkspaceResolver{
			Q:            r,
			WorkspaceRow: row,
		}
	}
	return resolvers, nil
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

func (r *MutationResolver) NewWorkspace(ctx context.Context, args struct {
	Root      string
	ProjectID *string
}) (*WorkspaceResolver, error) {
	var row WorkspaceRow
	row.ID = gensym.RandomBase32()
	row.Root = args.Root
	row.ProjectID = args.ProjectID
	if _, err := r.DB.ExecContext(ctx, `
		INSERT INTO workspace ( id, root, project_id )
		VALUES ( ?, ?, ? )
	`, row.ID, row.Root, row.ProjectID); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &WorkspaceResolver{
		Q:            r,
		WorkspaceRow: row,
	}, nil
}

func (r *WorkspaceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r.Q.projectByID(ctx, r.ProjectID)
}

func (r *WorkspaceResolver) StackID(ctx context.Context) (*string, error) {
	stack, err := r.Stack(ctx)
	if stack == nil || err != nil {
		return nil, err
	}
	return &stack.ID, nil
}

func (r *WorkspaceResolver) Stack(ctx context.Context) (*StackResolver, error) {
	stacks, err := r.Q.stacksByWorkspace(ctx, r.ID)
	if len(stacks) == 0 || err != nil {
		return nil, err
	}
	if len(stacks) > 1 {
		return nil, errors.New("ambiguous")
	}
	return stacks[0], nil
}
