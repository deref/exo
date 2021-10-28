package server

import (
	"context"
	"errors"
	"path"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/esv"
	"github.com/deref/exo/internal/manifest/exohcl"
)

func (ws *Workspace) getVaultURLs(ctx context.Context) []string {
	manifest := ws.tryLoadManifest(ctx)
	if manifest == nil {
		return nil
	}

	env := exohcl.NewEnvironment(manifest)
	_ = exohcl.Analyze(ctx, env)

	res := make([]string, 0, len(env.Secrets))
	for _, secrets := range env.Secrets {
		if secrets.Source == "" {
			continue
		}
		res = append(res, secrets.Source)
	}

	return res
}

func (ws *Workspace) AddVault(ctx context.Context, input *api.AddVaultInput) (*api.AddVaultOutput, error) {
	err := ws.modifyManifest(ctx, exohcl.AppendSecrets{
		Source: input.URL,
	})
	if err != nil {
		return nil, err
	}
	return &api.AddVaultOutput{}, nil
}

func (ws *Workspace) RemoveVault(ctx context.Context, input *api.RemoveVaultInput) (*api.RemoveVaultOutput, error) {
	err := ws.modifyManifest(ctx, exohcl.RemoveSecrets{
		Source: input.URL,
	})
	if err != nil {
		return nil, err
	}
	return &api.RemoveVaultOutput{}, nil
}

func (ws *Workspace) DescribeVaults(ctx context.Context, input *api.DescribeVaultsInput) (*api.DescribeVaultsOutput, error) {
	vaultURLs := ws.getVaultURLs(ctx)

	descriptions := make([]api.VaultDescription, len(vaultURLs))
	for i, vaultURL := range vaultURLs {
		// TODO: add a status command.
		_, err := ws.EsvClient.GetWorkspaceSecrets(vaultURL)
		descriptions[i] = api.VaultDescription{
			Name:      path.Base(vaultURL), // XXX
			URL:       vaultURL,
			Connected: err == nil,
			NeedsAuth: errors.Is(err, esv.AuthError),
		}
	}

	return &api.DescribeVaultsOutput{Vaults: descriptions}, nil
}
