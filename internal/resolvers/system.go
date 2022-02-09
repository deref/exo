package resolvers

import "context"

type SystemResolver struct {
	Q *RootResolver
}

func (r *RootResolver) System() *SystemResolver {
	return &SystemResolver{Q: r}
}

func (r *SystemResolver) Stream() *StreamResolver {
	return r.Q.streamForSource("System", "SYSTEM")
}

func (r *SystemResolver) eventPrototype(ctx context.Context) (EventRow, error) {
	return EventRow{
		SourceType: "System",
	}, nil
}
