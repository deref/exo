package resolvers

import "context"

type ClusterResolver struct {
	Q *QueryResolver
	ClusterRow
}

type ClusterRow struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (r *QueryResolver) ClusterByName(ctx context.Context, args struct {
	Name string
}) (*ClusterResolver, error) {
	return r.clusterByName(ctx, args.Name)
}

func (r *QueryResolver) clusterByName(ctx context.Context, name string) (*ClusterResolver, error) {
	clus := &ClusterResolver{}
	err := r.getRowByID(ctx, &clus.ClusterRow, `
		SELECT id, name
		FROM cluster
		WHERE name = ?
	`, &name)
	if clus.ID == "" {
		clus = nil
	}
	return clus, err
}

// NOTE [DEFAULT_CLUSTER]: The default cluster should be configurable, or at
// least optional.  Consider remote/CI use cases where no components/resources
// should be run locally.
const defaultClusterName = "local"

func (r *QueryResolver) DefaultCluster(ctx context.Context) (*ClusterResolver, error) {
	return r.clusterByName(ctx, "local")
}

// SEE NOTE [DEFAULT_CLUSTER].
func (r *ClusterResolver) Default() bool {
	return r.Name == "local"
}
