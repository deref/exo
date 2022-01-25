package resolvers

import (
	"context"
	"fmt"
)

type ClusterResolver struct {
	Q *QueryResolver
	ClusterRow
}

type ClusterRow struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (r *QueryResolver) clusterByID(ctx context.Context, id *string) (*ClusterResolver, error) {
	clus := &ClusterResolver{}
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
		return nil, fmt.Errorf("ambiguous cluster ref: %q", ref)
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