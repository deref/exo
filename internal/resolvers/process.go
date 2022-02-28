package resolvers

import "context"

type ProcessResolver struct {
	ID   string
	Name string
}

// TODO: Include running jobs?
func (r *QueryResolver) processesByStack(ctx context.Context, stackID string) ([]*ProcessResolver, error) {
	components, err := r.componentsByStack(ctx, stackID)
	if err != nil {
		return nil, err
	}

	processes := make([]*ProcessResolver, 0, len(components))
	for _, component := range components {
		if !component.isRunningProcess() {
			continue
		}
		processes = append(processes, &ProcessResolver{
			ID:   component.ID,
			Name: component.Name,
		})
	}
	return processes, nil
}
