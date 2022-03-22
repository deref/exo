package resolvers

import "context"

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
		process := r.processFromComponent(component)
		if process == nil {
			continue
		}
		// XXX Some callers may also want to see recently terminated processes.
		if !process.Running() {
			continue
		}
		processes = append(processes, process)
	}
	return processes, nil
}

func (r *ProcessComponentResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, &r.ComponentID)
}

func (r *ProcessResolver) CPUPercentage() *float64 {
	return nil // TODO!
}

func (r *ProcessResolver) Running() bool {
	return true // XXX
}
