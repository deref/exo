package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
	"github.com/jmoiron/sqlx"
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
		SELECT *
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
	err := r.getRowByKey(ctx, &s.StackRow, `
		SELECT *
		FROM stack
		WHERE id = ?
	`, id)
	if s.ID == "" {
		s = nil
	}
	return s, err
}

func (r *QueryResolver) stacksByWorkspaceID(ctx context.Context, workspaceID string) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT stack.*
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
	return r.stackByRef(ctx, &args.Ref)
}

func (r *QueryResolver) stackByRef(ctx context.Context, ref *string) (*StackResolver, error) {
	stack, err := r.stackByID(ctx, ref)
	if stack != nil || err != nil {
		return stack, err
	}
	ws, err := r.workspaceByRef(ctx, ref)
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
		ws, err = r.workspaceByRef(ctx, args.Workspace)
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

func (r *StackResolver) componentByRef(ctx context.Context, ref string) (*ComponentResolver, error) {
	return r.Q.componentByRef(ctx, ref, stringPtr(r.ID))
}

func (r *MutationResolver) SetWorkspaceStack(ctx context.Context, args struct {
	Workspace string
	Stack     *string
}) (*StackResolver, error) {
	workspace, err := r.workspaceByRef(ctx, &args.Workspace)
	if err != nil {
		return nil, fmt.Errorf("resolving workspace: %w", err)
	}
	if workspace == nil {
		return nil, fmt.Errorf("no such workspace: %q", args.Workspace)
	}
	var stackID *string
	if args.Stack != nil {
		stack, err := r.stackByRef(ctx, args.Stack)
		if err != nil {
			return nil, fmt.Errorf("resolving stack: %w", err)
		}
		if stack == nil {
			return nil, fmt.Errorf("no such stack: %q", *args.Stack)
		}
		stackID = &stack.ID
	}
	var stackRow StackRow
	err = transact(ctx, r.DB, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			UPDATE stack
			SET workspace_id = null
			WHERE workspace_id = ?;
		`, workspace.ID); err != nil {
			return err
		}
		return tx.GetContext(ctx, &stackRow, `
			UPDATE stack
			SET workspace_id = ?
			WHERE id = ?
			RETURNING *;
		`, workspace.ID, stackID)
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &StackResolver{
		Q:        r,
		StackRow: stackRow,
	}, nil
}
