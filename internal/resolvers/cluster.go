package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/manifest/exocue"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/shellutil"
)

type ClusterResolver struct {
	Q *RootResolver
	ClusterRow
}

type ClusterRow struct {
	ID                   string             `db:"id"`
	Name                 string             `db:"name"`
	EnvironmentVariables scalars.JSONObject `db:"environment_variables"`
	Updated              scalars.Instant    `db:"updated"`
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
	err := r.DB.SelectContext(ctx, &rows, `
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
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM cluster
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ClusterResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ClusterResolver{
			Q:          r,
			ClusterRow: row,
		}
	}
	return resolvers, nil
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
	Environment *scalars.JSONObject
}) (*ClusterResolver, error) {
	return r.updateCluster(ctx, args.Ref, args.Environment)
}

func (r *MutationResolver) updateCluster(ctx context.Context, ref string, environment *scalars.JSONObject) (*ClusterResolver, error) {
	now := scalars.Now(ctx)
	var row ClusterRow
	if err := r.DB.GetContext(ctx, &row, `
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

func (r *ClusterResolver) Environment(ctx context.Context) (*EnvironmentResolver, error) {
	variables := r.EnvironmentVariables

	now := scalars.Now(ctx)
	if r.EnvironmentVariables == nil || now.Sub(r.Updated) < clusterTTL {
		r, err := r.Q.refreshCluster(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		variables = r.EnvironmentVariables
	}

	return JSONObjectToEnvironment(variables, "Cluster")
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
	envObj := make(scalars.JSONObject, len(envMap))
	for k, v := range envMap {
		envObj[k] = v
	}
	return r.updateCluster(ctx, ref, &envObj)
}

func (r *ClusterResolver) Configuration(ctx context.Context) (string, error) {
	cfg, err := r.configuration(ctx)
	if err != nil {
		return "", err
	}
	return formatConfiguration(cue.Value(cfg))
}

func (r *ClusterResolver) configuration(ctx context.Context) (exocue.Cluster, error) {
	b := exocue.NewBuilder()
	if err := r.addConfiguration(ctx, b); err != nil {
		return exocue.Cluster{}, err
	}
	return b.BuildCluster(), nil
}

func (r *ClusterResolver) addConfiguration(ctx context.Context, b *exocue.Builder) error {
	env, err := r.Environment(ctx)
	if err != nil {
		return fmt.Errorf("resolving environment: %w", err)
	}
	b.AddCluster(r.ID, r.Name, env.AsMap())
	return nil
}
