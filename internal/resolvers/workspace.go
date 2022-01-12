package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/pathutil"
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

func (r *QueryResolver) WorkspaceByRef(ctx context.Context, args struct {
	Ref string
}) (*WorkspaceResolver, error) {
	workspaces, err := r.AllWorkspaces(ctx)
	if err != nil {
		return nil, err
	}
	var deepest *WorkspaceResolver
	maxLen := 0
	for _, workspace := range workspaces {
		// Exact match by ID.
		if workspace.ID == args.Ref {
			return workspace, nil
		}
		// Match by root. Prefer deepest root prefix match.
		n := len(workspace.Root)
		if n > maxLen && pathutil.HasFilePathPrefix(args.Ref, workspace.Root) {
			deepest = workspace
			maxLen = n
		}
	}
	return deepest, nil
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

func (r *WorkspaceResolver) Environment(ctx context.Context) *EnvironmentResolver {
	return &EnvironmentResolver{
		Workspace: r,
	}
}

func (r *WorkspaceResolver) FileSystem() *FileSystemResolver {
	return &FileSystemResolver{
		hostPath: r.Root,
	}
}
