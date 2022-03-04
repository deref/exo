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
	ID        string `db:"id"`
	Root      string `db:"root"`
	ProjectID string `db:"project_id"`
}

func workspaceRowsToResolvers(r *RootResolver, rows []WorkspaceRow) []*WorkspaceResolver {
	resolvers := make([]*WorkspaceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &WorkspaceResolver{
			Q:            r,
			WorkspaceRow: row,
		}
	}
	return resolvers
}

func (r *QueryResolver) AllWorkspaces(ctx context.Context) ([]*WorkspaceResolver, error) {
	var rows []WorkspaceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM workspace
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	return workspaceRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) WorkspaceByID(ctx context.Context, args struct {
	ID string
}) (*WorkspaceResolver, error) {
	return r.workspaceByID(ctx, &args.ID)
}

func (r *QueryResolver) workspaceByID(ctx context.Context, id *string) (*WorkspaceResolver, error) {
	ws := &WorkspaceResolver{}
	err := r.getRowByKey(ctx, &ws.WorkspaceRow, `
		SELECT *
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
	return r.workspaceByRef(ctx, &args.Ref)
}

func (r *QueryResolver) workspaceByRef(ctx context.Context, ref *string) (*WorkspaceResolver, error) {
	if ref == nil {
		return nil, nil
	}
	refStr := *ref
	workspaces, err := r.AllWorkspaces(ctx)
	if err != nil {
		return nil, err
	}
	var deepest *WorkspaceResolver
	maxLen := 0
	for _, workspace := range workspaces {
		// Exact match by ID.
		if workspace.ID == refStr {
			return workspace, nil
		}
		// Match by root. Prefer deepest root prefix match.
		n := len(workspace.Root)
		if n > maxLen && pathutil.HasFilePathPrefix(refStr, workspace.Root) {
			deepest = workspace
			maxLen = n
		}
	}
	return deepest, nil
}

func (r *MutationResolver) CreateWorkspace(ctx context.Context, args struct {
	Root      string
	ProjectID *string
}) (*WorkspaceResolver, error) {
	var row WorkspaceRow
	row.ID = gensym.RandomBase32()
	row.Root = args.Root

	if args.ProjectID == nil {
		proj, err := r.CreateProject(ctx, struct{ DisplayName *string }{})
		if err != nil {
			return nil, fmt.Errorf("creating new project for workspace: %w", err)
		}
		row.ProjectID = proj.ID
	} else {
		row.ProjectID = *args.ProjectID
	}

	if err := r.insertRow(ctx, "workspace", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &WorkspaceResolver{
		Q:            r,
		WorkspaceRow: row,
	}, nil
}

func (r *WorkspaceResolver) URL() string {
	return r.Q.Routes.workspaceURL(r.ID)
}

func (r *WorkspaceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r.Q.projectByID(ctx, &r.ProjectID)
}

func (r *WorkspaceResolver) StackID(ctx context.Context) (*string, error) {
	stack, err := r.Stack(ctx)
	if stack == nil || err != nil {
		return nil, err
	}
	return &stack.ID, nil
}

func (r *WorkspaceResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByWorkspaceID(ctx, r.ID)
}

func (r *WorkspaceResolver) componentByRef(ctx context.Context, ref string) (*ComponentResolver, error) {
	stack, err := r.Stack(ctx)
	if stack == nil || err != nil {
		return nil, err
	}
	return stack.componentByRef(ctx, ref)
}

func (r *WorkspaceResolver) Environment(ctx context.Context) (*EnvironmentResolver, error) {
	stack, err := r.Stack(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving stack: %w", err)
	}
	if stack == nil {
		// TODO: When there is no current stack, it may be useful to return the
		// environment of the workspace's cluster/host.
		return nil, errors.New("no current stack")
	}
	return stack.Environment(ctx)
}

func (r *WorkspaceResolver) FileSystem() *FileSystemResolver {
	return &FileSystemResolver{
		hostPath: r.Root,
	}
}

func (r *MutationResolver) BuildWorkspace(ctx context.Context, args struct {
	Workspace string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) BuildWorkspaceComponents(ctx context.Context, args struct {
	Workspace  string
	Components []string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) RefreshWorkspace(ctx context.Context, args struct {
	Workspace string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) RefreshWorkspaceComponents(ctx context.Context, args struct {
	Workspace  string
	Components []string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) StartWorkspace(ctx context.Context, args struct {
	Workspace string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) StartWorkspaceComponents(ctx context.Context, args struct {
	Workspace  string
	Components []string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) RestartWorkspace(ctx context.Context, args struct {
	Workspace string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) RestartWorkspaceComponents(ctx context.Context, args struct {
	Workspace  string
	Components []string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) StopWorkspace(ctx context.Context, args struct {
	Workspace string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *MutationResolver) StopWorkspaceComponents(ctx context.Context, args struct {
	Workspace  string
	Components []string
}) (*TaskResolver, error) {
	return nil, errors.New("NOT YET IMPLEMENTED")
}

func (r *WorkspaceResolver) Components(ctx context.Context, args struct {
	All       *bool
	Recursive *bool
}) ([]*ComponentResolver, error) {
	stack, err := r.Stack(ctx)
	if stack == nil || err != nil {
		return nil, err
	}
	return stack.Components(ctx, struct {
		All       *bool
		Recursive *bool
	}{
		All:       args.All,
		Recursive: args.Recursive,
	})
}

func (r *WorkspaceResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	stack, err := r.Stack(ctx)
	if stack == nil || err != nil {
		return nil, err
	}
	return stack.Resources(ctx)
}

func (r *WorkspaceResolver) Manifest(ctx context.Context) (*ManifestResolver, error) {
	// XXX Use active manifest from workspace configuration.
	var format *string
	return r.findManifest(ctx, format)
}

func (r *WorkspaceResolver) FindManifest(ctx context.Context, args struct {
	Format *string
}) (*ManifestResolver, error) {
	return r.findManifest(ctx, args.Format)
}

func (r *WorkspaceResolver) findManifest(ctx context.Context, format *string) (*ManifestResolver, error) {
	return r.Q.findManifest(ctx, r.FileSystem(), format)
}

func (r *WorkspaceResolver) manifestByPath(ctx context.Context, path string, format *string) (*ManifestResolver, error) {
	return r.Q.manifestByPath(ctx, r.FileSystem(), path, format)
}
