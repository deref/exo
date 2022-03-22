package resolvers

import (
	"context"
)

type StoreResolver struct {
	Q    *RootResolver
	Type string
}

type StoreComponentResolver struct {
	StoreResolver

	ComponentID string
	Name        string
}

func (r *QueryResolver) isStoreType(typ string) bool {
	// TODO: Extensible.
	switch typ {
	case "volume":
		return true
	default:
		return false
	}
}

func (r *QueryResolver) storeFromComponent(component *ComponentResolver) *StoreComponentResolver {
	if !r.isStoreType(component.Type) {
		return nil
	}
	return &StoreComponentResolver{
		StoreResolver: StoreResolver{
			Q:    r,
			Type: component.Type,
		},
		ComponentID: component.ID,
		Name:        component.Name,
	}
}

func (r *QueryResolver) storesByStack(ctx context.Context, stackID string) ([]*StoreComponentResolver, error) {
	components, err := r.componentsByStack(ctx, stackID)
	if err != nil {
		return nil, err
	}

	stores := make([]*StoreComponentResolver, 0, len(components))
	for _, component := range components {
		store := r.storeFromComponent(component)
		if store == nil {
			continue
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (r *StoreComponentResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, &r.ComponentID)
}

func (r *StoreResolver) SizeMiB() *float64 {
	return nil // TODO!
}
