package resolvers

import (
	"context"
	"time"
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

func (r *SystemResolver) Version() *VersionInfoResolver {
	return resolveVersionInfo()
}

func (r *RootResolver) SystemChange(ctx context.Context) <-chan *SystemResolver {
	c := make(chan *SystemResolver, 1)

	var upgrade *string
	{
		system := r.System()
		c <- system
		upgrade = system.Version().Upgrade(ctx)
	}

	// Poll for available upgrades.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(60 * time.Second):
				system := r.System()
				newUpgrade := system.Version().Upgrade(ctx)
				if newUpgrade != nil && (upgrade == nil || *upgrade != *newUpgrade) {
					upgrade = newUpgrade
					c <- system
				}
			}
		}
	}()

	return c
}
