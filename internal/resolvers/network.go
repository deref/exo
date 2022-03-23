package resolvers

import "context"

type NetworkResolver struct {
	Q    *RootResolver
	Type string
}

type NetworkComponentResolver struct {
	NetworkResolver

	ComponentID string
	Name        string
}

func (r *QueryResolver) isNetworkType(typ string) bool {
	// TODO: Extensible.
	switch typ {
	case "network":
		return true
	default:
		return false
	}
}

func (r *QueryResolver) networkFromComponent(component *ComponentResolver) *NetworkComponentResolver {
	if !r.isNetworkType(component.Type) {
		return nil
	}
	return &NetworkComponentResolver{
		NetworkResolver: NetworkResolver{
			Q:    r,
			Type: component.Type,
		},
		ComponentID: component.ID,
		Name:        component.Name,
	}
}

func (r *QueryResolver) networksByStack(ctx context.Context, stackID string) ([]*NetworkComponentResolver, error) {
	components, err := r.componentsByStack(ctx, stackID)
	if err != nil {
		return nil, err
	}

	networks := make([]*NetworkComponentResolver, 0, len(components))
	for _, component := range components {
		network := r.networkFromComponent(component)
		if network == nil {
			continue
		}
		networks = append(networks, network)
	}
	return networks, nil
}

func (r *NetworkComponentResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, &r.ComponentID)
}
