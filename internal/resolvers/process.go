package resolvers

import (
	"context"

	. "github.com/deref/exo/internal/scalars"
)

type ProcessResolver struct {
	Q    *RootResolver
	Type string
}

type ProcessComponentResolver struct {
	ProcessResolver

	ComponentID string
	Name        string
}

func (r *QueryResolver) isProcessType(typ string) bool {
	// TODO: Extensible.
	switch typ {
	case "daemon", "process", "container":
		return true
	default:
		return false
	}
}

func (r *QueryResolver) processFromComponent(component *ComponentResolver) *ProcessComponentResolver {
	if r.isProcessType(component.Type) {
		return nil
	}
	return &ProcessComponentResolver{
		ProcessResolver: ProcessResolver{
			Q:    r,
			Type: component.Type,
		},
		ComponentID: component.ID,
		Name:        component.Name,
	}
}

// TODO: Include running jobs?
func (r *QueryResolver) processesByStack(ctx context.Context, stackID string) ([]*ProcessComponentResolver, error) {
	components, err := r.componentsByStack(ctx, stackID)
	if err != nil {
		return nil, err
	}

	processes := make([]*ProcessComponentResolver, 0, len(components))
	for _, component := range components {
		// TODO: Some callers may also want to see recently terminated processes.
		if !component.Running() {
			continue
		}
		process := r.processFromComponent(component)
		if process == nil {
			continue
		}
		processes = append(processes, process)
	}
	return processes, nil
}

func (r *ProcessComponentResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, &r.ComponentID)
}

func (r *ProcessResolver) Started() *Instant {
	return nil // TODO!
}

func (r *ProcessResolver) CPUPercent() *float64 {
	return nil // TODO!
}

func (r *ProcessResolver) ResidentBytes() *int32 {
	return nil // TODO!
}

func (r *ProcessResolver) Ports() *[]int32 {
	return nil // TODO!
}

func (r *ProcessResolver) Environment() *EnvironmentResolver {
	return nil // TODO!
}
