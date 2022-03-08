package resolvers

import (
	"context"

	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/logging"
)

type VersionInfoResolver struct {
	Installed string
	Build     string
	Managed   bool
}

func resolveVersionInfo() *VersionInfoResolver {
	return &VersionInfoResolver{
		Installed: about.Version,
		Build:     about.GetBuild(),
		Managed:   about.IsManaged,
	}
}

func (r *VersionInfoResolver) Latest(ctx context.Context) (string, error) {
	tel := telemetry.FromContext(ctx)
	return tel.LatestVersion(ctx)
}

func (r *VersionInfoResolver) Upgrade(ctx context.Context) *string {
	latest, err := r.Latest(ctx)
	if err != nil {
		logging.Infof(ctx, "resolving latest version: %w", err)
		return nil
	}
	if latest <= r.Installed {
		return nil
	}
	return &latest
}
