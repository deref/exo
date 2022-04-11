package resolvers

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/manifest/exocue"
)

type ConfigurationResolver struct {
	Q         *QueryResolver
	StackID   string
	Recursive bool
	Final     bool
}

func (r *QueryResolver) fullConfiguration(stackID string) *ConfigurationResolver {
	return &ConfigurationResolver{
		Q:         r,
		StackID:   stackID,
		Recursive: true,
		Final:     true,
	}
}

func (r *QueryResolver) fullConfigurationAsCueValue(ctx context.Context, stackID string) (exocue.Configuration, error) {
	return r.fullConfiguration(stackID).AsCueValue(ctx)
}

func (r *ConfigurationResolver) StackAsString(ctx context.Context) (string, error) {
	built, err := r.AsCueValue(ctx)
	if err != nil {
		return "", err
	}
	return r.valueAsString(cue.Value(built.Stack()))
}

func (r *ConfigurationResolver) ComponentAsString(ctx context.Context, componentID string) (string, error) {
	built, err := r.AsCueValue(ctx)
	if err != nil {
		return "", err
	}
	return r.valueAsString(cue.Value(built.Component(componentID)))
}

func (r *ConfigurationResolver) valueAsString(v cue.Value) (string, error) {
	return formatConfiguration(v, r.Final)
}

// TODO: How much does it help if this is cached?
func (r *ConfigurationResolver) AsCueValue(ctx context.Context) (exocue.Configuration, error) {
	b := exocue.NewBuilder()

	stack, err := r.Q.stackByID(ctx, &r.StackID)
	if err := validateResolve("stack", r.StackID, stack, err); err != nil {
		return exocue.Configuration{}, err
	}
	b.SetStack(stack.ID, stack.Name)

	cluster, err := stack.Cluster(ctx)
	if err := validateResolve("cluster", stack.ClusterID, cluster, err); err != nil {
		return exocue.Configuration{}, err
	}
	clusterEnvironment, err := cluster.Environment(ctx)
	if err != nil {
		return exocue.Configuration{}, fmt.Errorf("resolving cluster environment: %w", err)
	}
	b.SetCluster(cluster.ID, cluster.Name, clusterEnvironment.asMap())

	componentSet := &componentSetResolver{
		Q:         r.Q,
		StackID:   r.StackID,
		Recursive: r.Recursive,
	}
	components, err := componentSet.Items(ctx)
	if err != nil {
		return exocue.Configuration{}, fmt.Errorf("resolving components: %w", err)
	}
	for _, component := range components {
		b.AddComponent(component.ID, component.Name, component.Type, cue.Value(component.Spec), component.ParentID)
	}

	resources, err := stack.Resources(ctx)
	if err != nil {
		return exocue.Configuration{}, fmt.Errorf("resolving resources: %w", err)
	}
	for _, resource := range resources {
		b.AddResource(resource.ID, resource.Type, resource.IRI, resource.ComponentID)
	}

	return b.Build()
}
