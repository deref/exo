package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/deref/exo/internal/gensym"
	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/jmoiron/sqlx"
)

type StackResolver struct {
	Q *QueryResolver
	StackRow
}

type StackRow struct {
	ID          string   `db:"id"`
	Name        string   `db:"name"`
	ClusterID   string   `db:"cluster_id"`
	ProjectID   *string  `db:"project_id"`
	WorkspaceID *string  `db:"workspace_id"`
	Disposed    *Instant `db:"disposed"`
}

func stackRowsToResolvers(r *RootResolver, rows []StackRow) []*StackResolver {
	resolvers := make([]*StackResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &StackResolver{
			Q:        r,
			StackRow: row,
		}
	}
	return resolvers
}

func (r *QueryResolver) AllStacks(ctx context.Context) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM stack
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	return stackRowsToResolvers(r, rows), nil
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
	err := r.db.SelectContext(ctx, &rows, `
		SELECT stack.id, stack.name, stack.cluster_id, stack.project_id, workspace_id
		FROM stack
		INNER JOIN project ON stack.project_id = project.id
		WHERE project.id = ?
		ORDER BY stack.name ASC
	`, stackID)
	if err != nil {
		return nil, err
	}
	return stackRowsToResolvers(r, rows), nil
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

func (r *StackResolver) DisplayName() string {
	return r.Name // TODO: This is likely to be ambiguous, so use workspace display name somehow.
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
	componentSet.All = isTrue(args.All)
	componentSet.Recursive = isTrue(args.Recursive)
	return componentSet.Items(ctx)
}

func (r *StackResolver) components(ctx context.Context) ([]*ComponentResolver, error) {
	return r.Q.componentsByStack(ctx, r.ID)
}

func (r *StackResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByStack(ctx, r.ID)
}

func (r *StackResolver) Processes(ctx context.Context) ([]*ProcessComponentResolver, error) {
	return r.Q.processesByStack(ctx, r.ID)
}

func (r *StackResolver) Stores(ctx context.Context) ([]*StoreComponentResolver, error) {
	return r.Q.storesByStack(ctx, r.ID)
}

func (r *StackResolver) Networks(ctx context.Context) ([]*NetworkComponentResolver, error) {
	return r.Q.networksByStack(ctx, r.ID)
}

func (r *MutationResolver) CreateStack(ctx context.Context, args struct {
	Workspace *string
	Name      *string
	Cluster   *string
}) (*StackResolver, error) {
	var workspace *WorkspaceResolver
	if args.Workspace != nil {
		var err error
		workspace, err = r.workspaceByRef(ctx, args.Workspace)
		if err := validateResolve("workspace", *args.Workspace, workspace, err); err != nil {
			return nil, err
		}
	}

	var cluster *ClusterResolver
	if args.Cluster == nil {
		var err error
		cluster, err = r.DefaultCluster(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving default cluster: %q", err)
		}
		if cluster == nil {
			return nil, errutil.HTTPErrorf(http.StatusNotFound, "no default cluster")
		}
	} else {
		var err error
		cluster, err = r.clusterByRef(ctx, *args.Cluster)
		if err := validateResolve("cluster", *args.Cluster, cluster, err); err != nil {
			return nil, err
		}
	}

	var row StackRow
	row.ID = gensym.RandomBase32()
	row.Name = *trimmedPtr(args.Name, row.ID)
	row.ClusterID = cluster.ID
	if workspace != nil {
		row.WorkspaceID = &workspace.ID
		row.ProjectID = &workspace.ProjectID
	}

	// TODO: Validate name.

	if _, err := r.db.ExecContext(ctx, `
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

func (r *MutationResolver) DestroyStack(ctx context.Context, args struct {
	Ref string
}) (*ReconciliationResolver, error) {
	stack, err := r.stackByRef(ctx, &args.Ref)
	if err := validateResolve("stack", args.Ref, stack, err); err != nil {
		return nil, err
	}
	stack, err = r.disposeStack(ctx, stack.ID)
	if err != nil {
		return nil, err
	}
	return r.startStackReconciliation(ctx, stack)
}

func (r *MutationResolver) disposeStack(ctx context.Context, id string) (*StackResolver, error) {
	now := Now(ctx)
	var row StackRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE stack
		SET disposed = COALESCE(disposed, ?)
		WHERE id = ?
		RETURNING *
	`, now, id,
	); err != nil {
		return nil, err
	}
	if err := r.disposeComponentsByStack(ctx, id); err != nil {
		return nil, fmt.Errorf("disposing stack components: %w", err)
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
	if err := validateResolve("workspace", args.Workspace, workspace, err); err != nil {
		return nil, err
	}
	var stackID *string
	if args.Stack != nil {
		stack, err := r.stackByRef(ctx, args.Stack)
		if err := validateResolve("stack", *args.Stack, stack, err); err != nil {
			return nil, err
		}
		stackID = &stack.ID
	}
	var stackRow StackRow
	err = transact(ctx, r.db, func(tx *sqlx.Tx) error {
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

func (r *StackResolver) Configuration(ctx context.Context, args struct {
	Recursive *bool
	Final     *bool
}) (string, error) {
	configuration := &ConfigurationResolver{
		Q:         r.Q,
		StackID:   r.ID,
		Recursive: isTrue(args.Recursive),
		Final:     isTrue(args.Final),
	}
	return configuration.ComponentAsString(ctx, r.ID)
}

func (r *StackResolver) Environment(ctx context.Context) (*EnvironmentResolver, error) {
	// XXX figure me out!
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

func (r *StackResolver) Vaults(ctx context.Context) ([]*VaultResolver, error) {
	return r.Q.vaultsByStackID(ctx, r.ID)
}

func (r *StackResolver) Secrets(ctx context.Context) ([]*SecretResolver, error) {
	return r.Q.secretsByStackID(ctx, r.ID)
}
