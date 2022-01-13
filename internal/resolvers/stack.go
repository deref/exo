package resolvers

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type StackResolver struct {
	Q *QueryResolver
	StackRow
}

type StackRow struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	ClusterID   string  `db:"cluster_id"`
	ProjectID   *string `db:"project_id"`
	WorkspaceID *string `db:"workspace_id"`
}

func (r *QueryResolver) AllStacks(ctx context.Context) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT id, name, cluster_id, project_id, workspace_id
		FROM stack
		ORDER BY id ASC
	`)
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

func (r *QueryResolver) StackByID(ctx context.Context, args struct {
	ID string
}) (*StackResolver, error) {
	return r.stackByID(ctx, &args.ID)
}

func (r *QueryResolver) stackByID(ctx context.Context, id *string) (*StackResolver, error) {
	s := &StackResolver{}
	err := r.getRowByID(ctx, &s.StackRow, `
		SELECT id, name, cluster_id, project_id, workspace_id
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
		SELECT stack.id, stack.name, stack.cluster_id, stack.project_id, workspace_id
		FROM stack
		INNER JOIN workspace ON stack.workspace_id = workspace.id
		WHERE workspace.id = ?
		ORDER BY stack.name ASC
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

func (r *QueryResolver) stacksByProject(ctx context.Context, stackID string) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT stack.id, stack.name, stack.cluster_id, stack.project_id, workspace_id
		FROM stack
		INNER JOIN project ON stack.project_id = project.id
		WHERE project.id = ?
		ORDER BY stack.name ASC
	`, stackID)
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

func (r *StackResolver) Cluster(ctx context.Context) (*ClusterResolver, error) {
	return r.Q.clusterByID(ctx, &r.ClusterID)
}

func (r *StackResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r.Q.projectByID(ctx, r.ProjectID)
}

func (r *StackResolver) Workspace(ctx context.Context) (*WorkspaceResolver, error) {
	return r.Q.workspaceByID(ctx, r.WorkspaceID)
}

func (r *StackResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByStack(ctx, r.ID)
}

func (r *MutationResolver) NewStack(ctx context.Context, args struct {
	Workspace *string
	Name      *string
	Cluster   *string
}) (*StackResolver, error) {
	var ws *WorkspaceResolver
	if args.Workspace != nil {
		var err error
		ws, err = r.workspaceByRef(ctx, *args.Workspace)
		if err != nil {
			return nil, fmt.Errorf("resolving workspace ref: %w", err)
		}
		if ws == nil {
			return nil, fmt.Errorf("no such workspace: %q", *args.Workspace)
		}
	}

	var clus *ClusterResolver
	if args.Cluster == nil {
		var err error
		clus, err = r.DefaultCluster(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving default cluster: %q", err)
		}
		if clus == nil {
			return nil, fmt.Errorf("no default cluster")
		}
	} else {
		var err error
		clus, err = r.clusterByRef(ctx, *args.Cluster)
		if err != nil {
			return nil, fmt.Errorf("resolving cluster: %q", err)
		}
		if clus == nil {
			return nil, fmt.Errorf("no such cluster: %q", *args.Cluster)
		}
	}

	var row StackRow
	row.ID = gensym.RandomBase32()
	row.Name = *trimmedPtr(args.Name, row.ID)
	row.ClusterID = clus.ID
	if ws != nil {
		row.WorkspaceID = &ws.ID
		row.ProjectID = &ws.ProjectID
	}

	// TODO: Validate name.

	if _, err := r.DB.ExecContext(ctx, `
		BEGIN;

		UPDATE stack
		SET workspace_id = NULL
		WHERE workspace_id = ?;

		INSERT INTO stack ( id, name, cluster_id, project_id, workspace_id )
		VALUES ( ?, ?, ?, ?, ? );

		COMMIT
	`, row.WorkspaceID, row.ID, row.Name, row.ClusterID, row.ProjectID, row.WorkspaceID); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &StackResolver{
		Q:        r,
		StackRow: row,
	}, nil
}
