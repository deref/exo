package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/manifest/exocue"
	"github.com/deref/exo/internal/util/errutil"
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
	stack := &StackResolver{
		Q: r,
	}
	err := r.getRowByKey(ctx, &stack.StackRow, `
		SELECT *
		FROM stack
		WHERE id = ?
	`, id)
	if stack.ID == "" {
		stack = nil
	}
	return stack, err
}

func (r *QueryResolver) stackByWorkspaceID(ctx context.Context, workspaceID string) (*StackResolver, error) {
	stack := &StackResolver{
		Q: r,
	}
	err := r.getRowByKey(ctx, &stack.StackRow, `
		SELECT stack.*
		FROM stack
		INNER JOIN workspace ON stack.workspace_id = workspace.id
		WHERE workspace.id = ?
		ORDER BY stack.name ASC
	`, &workspaceID)
	if stack.ID == "" {
		stack = nil
	}
	return stack, err
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

func (r *QueryResolver) stackByProjectIDAndRef(ctx context.Context, projectID string, ref string) (*StackResolver, error) {
	// Could move the filtering to the db-side, but not a big deal.
	stack, err := r.stackByRef(ctx, &ref)
	if stack != nil && (stack.ProjectID == nil || *stack.ProjectID != projectID) {
		stack = nil
	}
	return stack, err
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

func (r *StackResolver) Components(ctx context.Context, args struct {
	All       *bool
	Recursive *bool
}) ([]*ComponentResolver, error) {
	componentSet := &componentSetResolver{
		Q:       r.Q,
		StackID: r.ID,
	}
	if args.All != nil {
		componentSet.All = *args.All
	}
	if args.Recursive != nil {
		componentSet.Recursive = *args.Recursive
	}
	return componentSet.items(ctx)
}

func (r *StackResolver) components(ctx context.Context) ([]*ComponentResolver, error) {
	return r.Q.componentsByStack(ctx, r.ID)
}

func (r *StackResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByStack(ctx, r.ID)
}

func (r *StackResolver) Processes(ctx context.Context) ([]*ProcessResolver, error) {
	return r.Q.processesByStack(ctx, r.ID)
}

func (r *MutationResolver) CreateStack(ctx context.Context, args struct {
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

func (r *MutationResolver) RefreshStack(ctx context.Context, args struct {
	Ref string
}) (*ReconciliationResolver, error) {
	return nil, errors.New("TODO: refreshStack")
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
			return nil, errutil.HTTPErrorf(http.StatusNotFound, "no such stack: %q", *args.Stack)
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

func (r *StackResolver) Configuration(ctx context.Context) (string, error) {
	cfg, err := r.configuration(ctx)
	if err != nil {
		return "", err
	}
	return exocue.FormatString(cfg.Final())
}

// TODO: It might be valuable to cache this for multiple
// ComponentResolver.evalSpec calls.
func (r *StackResolver) configuration(ctx context.Context) (exocue.Stack, error) {
	b := exocue.NewBuilder()
	if err := r.addConfiguration(ctx, b); err != nil {
		return exocue.Stack{}, err
	}
	return b.BuildStack(), nil
}

func (r *StackResolver) addConfiguration(ctx context.Context, b *exocue.Builder) error {
	cluster, err := r.Cluster(ctx)
	if err != nil {
		return fmt.Errorf("resolving cluster: %w", err)
	}
	cluster.addConfiguration(ctx, b)

	components, err := r.components(ctx)
	if err != nil {
		return fmt.Errorf("resolving components: %w", err)
	}
	for _, component := range components {
		b.AddComponent(component.ID, component.Name, component.Type, component.Spec)
	}

	resources, err := r.Resources(ctx)
	if err != nil {
		return fmt.Errorf("resolving resources: %w", err)
	}
	for _, resource := range resources {
		b.AddResource(resource.ID, resource.Type, resource.IRI, resource.ComponentID)
	}
	return nil
}

func (r *StackResolver) Environment(ctx context.Context) (*EnvironmentResolver, error) {
	// XXX implement me
	// XXX This now does network requests and non-trivial parsing work. Therefore,
	// it is no longer appropriate to call deep in the call stack.
	/*
		ws := r.Workspace

		var sources []environment.Source

		if manifest := ws.tryLoadManifest(ctx); manifest != nil {
			manifestEnv := &exohcl.Environment
				Blocks: manifest.Environment,
			}
			diags := exohcl.Analyze(ctx, manifestEnv)
			if diags.HasErrors() {
				return nil, diags
			}
			sources = append(sources, manifestEnv)
		}

		sources = append(sources,
			environment.Default,
			&environment.OS{},
		)

		envPath, err := ws.resolveWorkspacePath(ctx, ".env")
		if err != nil {
			return nil, fmt.Errorf("resolving env file path: %w", err)
		}
		if exists, _ := osutil.Exists(envPath); exists {
			sources = append(sources, &environment.Dotenv{
				Path: envPath,
			})
		}

		b := &environmentBuilder{
			Environment: make(map[string]api.VariableDescription),
		}

		// TODO: Do not use DescribeVaults, instead build up sources from the
		// environment blocks ASTs. For example, there maybe a `variables` block or
		// some other environment sources that are not in the DescribeVaults output.
		describeVaultsResult, err := ws.DescribeVaults(ctx, &api.DescribeVaultsInput{})
		if err != nil {
			return nil, fmt.Errorf("getting vaults: %w", err)
		}

		logger := logging.CurrentLogger(ctx)
		for _, vault := range describeVaultsResult.Vaults {
			derefSource := &environment.ESV{
				Client: ws.EsvClient,
				Name:   vault.URL, // XXX
				URL:    vault.URL,
			}
			if err := derefSource.ExtendEnvironment(b); err != nil {
				// It's not appropriate to fail on error since this error could just
				// indicate the user is offline and thus cannot retrieve this value from
				// the secret provider.
				// TODO: this should really alert the user in a more apparent way that
				// fetching secrets from the vault has failed.
				logger.Infof("Could not extend environment from vault %q: %v", vault.URL, err)
			}
		}

		for _, source := range sources {
			if err := source.ExtendEnvironment(b); err != nil {
				return nil, fmt.Errorf("extending environment from %s: %w", source.EnvironmentSource(), err)
			}
		}
		return b.Environment, nil
	*/

	return &EnvironmentResolver{
		Variables: []*VariableResolver{
			{
				Name:  "X",
				Value: "a\b",
			},
		}, // XXX
	}, nil
}
