package resolvers

import (
	"context"

	"github.com/deref/exo/internal/about"
)

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

func (r *SystemResolver) Version() string {
	return about.Version
}

func (r *SystemResolver) Build() string {
	return about.GetBuild()
}
