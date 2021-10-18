package server

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/environment"
	"github.com/deref/exo/internal/util/osutil"
)

// XXX This now does network requests and non-trivial parsing work. Therefore,
// it is no longer appropriate to call deep in the call stack.
func (ws *Workspace) getEnvironment(ctx context.Context) (map[string]api.VariableDescription, error) {
	var sources []environment.Source

	describeVaultsResult, err := ws.DescribeVaults(ctx, &api.DescribeVaultsInput{})
	if err != nil {
		return nil, fmt.Errorf("getting vaults: %w", err)
	}
	for _, vault := range describeVaultsResult.Vaults {
		sources = append(sources, &environment.ESV{
			Client: ws.EsvClient,
			Name:   vault.Name,
			URL:    vault.Url,
		})
	}

	if manifest := ws.tryLoadManifest(ctx); manifest != nil {
		sources = append(sources, manifest.Environment())
	}

	sources = append(sources,
		environment.Default,
		&environment.OS{},
	)

	envPath, err := ws.resolveWorkspacePath(ctx, ".env")
	if err != nil {
		return nil, fmt.Errorf("resolving env file path: %w", err)
	}
	if exists, _ := osutil.Exists(envPath); exists {
		sources = append(sources, &environment.Dotenv{
			Path: envPath,
		})
	}

	b := &environmentBuilder{
		Environment: make(map[string]api.VariableDescription),
	}
	for _, source := range sources {
		if err := source.ExtendEnvironment(b); err != nil {
			return nil, fmt.Errorf("extending environment from %s: %w", source.EnvironmentSource(), err)
		}
	}
	return b.Environment, nil
}

type environmentBuilder struct {
	Environment map[string]api.VariableDescription
}

func (b *environmentBuilder) AppendVariable(src environment.Source, name string, value string) {
	b.Environment[name] = api.VariableDescription{
		Value:  value,
		Source: src.EnvironmentSource(),
	}
}
