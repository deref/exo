package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/shellutil"
)

type ClusterResolver struct {
	Q *RootResolver
	ClusterRow
}

type ClusterRow struct {
	ID                   string     `db:"id"`
	Name                 string     `db:"name"`
	EnvironmentVariables JSONObject `db:"environment_variables"`
	Updated              Instant    `db:"updated"`
}

func (r *QueryResolver) clusterByID(ctx context.Context, id *string) (*ClusterResolver, error) {
	clus := &ClusterResolver{
		Q: r,
	}
	err := r.getRowByKey(ctx, &clus.ClusterRow, `
		SELECT *
		FROM cluster
		WHERE id = ?
	`, id)
	if clus.ID == "" {
		clus = nil
	}
	return clus, err
}

func (r *QueryResolver) ClusterByRef(ctx context.Context, args struct {
	Ref string
}) (*ClusterResolver, error) {
	return r.clusterByRef(ctx, args.Ref)
}

func (r *QueryResolver) clusterByRef(ctx context.Context, ref string) (*ClusterResolver, error) {
	var rows []ClusterRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM cluster
		WHERE id = ? OR name = ?
	`, ref, ref)
	if len(rows) == 0 || err != nil {
		return nil, err
	}
	if len(rows) > 1 {
		return nil, errutil.HTTPErrorf(http.StatusConflict, "ambiguous cluster ref: %q", ref)
	}
	return &ClusterResolver{
		Q:          r,
		ClusterRow: rows[0],
	}, err
}

func (r *QueryResolver) AllClusters(ctx context.Context) ([]*ClusterResolver, error) {
	var rows []ClusterRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM cluster
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	return clusterRowsToResolvers(r, rows), nil
}

func clusterRowsToResolvers(r *RootResolver, rows []ClusterRow) []*ClusterResolver {
	resolvers := make([]*ClusterResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ClusterResolver{
			Q:          r,
			ClusterRow: row,
		}
	}
	return resolvers
}

// NOTE [DEFAULT_CLUSTER]: The default cluster should be configurable, or at
// least optional.  Consider remote/CI use cases where no components/resources
// should be run locally.
const defaultClusterName = "local"

func (r *QueryResolver) DefaultCluster(ctx context.Context) (*ClusterResolver, error) {
	return r.clusterByRef(ctx, "local")
}

// SEE NOTE [DEFAULT_CLUSTER].
func (r *ClusterResolver) Default() bool {
	return r.Name == "local"
}

func (r *MutationResolver) UpdateCluster(ctx context.Context, args struct {
	Ref         string
	Environment *JSONObject
}) (*ClusterResolver, error) {
	return r.updateCluster(ctx, args.Ref, args.Environment)
}

func (r *MutationResolver) updateCluster(ctx context.Context, ref string, environment *JSONObject) (*ClusterResolver, error) {
	now := Now(ctx)
	var row ClusterRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE cluster
		SET
			environment_variables = COALESCE(?, environment_variables),
			updated = ?
		WHERE id = ? OR name = ?
		RETURNING *
	`,
		environment,
		now,
		ref, ref,
	); err != nil {
		return nil, err
	}
	return &ClusterResolver{
		Q:          r,
		ClusterRow: row,
	}, nil
}

const clusterTTL = 3 * time.Second

// Cluster environments are not influenced by manifest files and so can be resolved directly.
// XXX If the environment is observed to change, should that trigger a reconciliation?
// XXX Alternatively, should it alert the user and allow them to take some action?
func (r *ClusterResolver) Environment(ctx context.Context) (*EnvironmentResolver, error) {
	locals := r.EnvironmentVariables

	now := Now(ctx)
	if r.EnvironmentVariables == nil || now.Sub(r.Updated) < clusterTTL {
		r, err := r.Q.refreshCluster(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		locals = r.EnvironmentVariables
	}

	environment := &EnvironmentResolver{
		Parent: nil,
		Source: r,
	}
	environment.initLocalsFromJSONObject(locals)
	return environment, nil
}

func (r *MutationResolver) RefreshCluster(ctx context.Context, args struct {
	Ref string
}) (*ClusterResolver, error) {
	return r.refreshCluster(ctx, args.Ref)
}

func (r *MutationResolver) refreshCluster(ctx context.Context, ref string) (*ClusterResolver, error) {
	envMap, err := shellutil.GetUserEnvironment(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying environment: %w", err)
	}
	envObj := make(JSONObject, len(envMap))
	for k, v := range envMap {
		envObj[k] = v
	}
	return r.updateCluster(ctx, ref, &envObj)
}
