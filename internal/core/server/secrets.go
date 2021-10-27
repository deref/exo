package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/esv"
)

var secretsUrlFile = "exo-secrets-url"

type vaultConfig struct {
	name string
	url  string
}

func (ws *Workspace) getVaultConfigs(ctx context.Context) ([]vaultConfig, error) {
	// NOTE: [VAULTS_IN_MANIFEST] Right now getVaultConfigs relies on an
	// exo-secrets-url file in the workspace root to provide the vault. At the
	// moment only one is supported. Instead we should add this to the state of
	// the workspace and support multiple vaults.

	secretConfigPath, err := ws.resolveWorkspacePath(ctx, secretsUrlFile)
	if err != nil {
		return nil, fmt.Errorf("resolving secrets config file path: %w", err)
	}
	secretsUrlBytes, err := ioutil.ReadFile(secretConfigPath)
	if os.IsNotExist(err) {
		return []vaultConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading secrets config: %w", err)
	}
	secretsUrl := string(bytes.TrimSpace(secretsUrlBytes))
	return []vaultConfig{{name: "esv-vault", url: secretsUrl}}, nil
}

// AddVault performs an upsert for the specified vault name and URL.
func (ws *Workspace) AddVault(ctx context.Context, input *api.AddVaultInput) (*api.AddVaultOutput, error) {
	// SEE NOTE [VAULTS_IN_MANIFEST]
	secretConfigPath, err := ws.resolveWorkspacePath(ctx, secretsUrlFile)
	if err != nil {
		return nil, fmt.Errorf("resolving secrets config file path: %w", err)
	}

	if err := ioutil.WriteFile(secretConfigPath, []byte(input.URL), 0600); err != nil {
		return nil, fmt.Errorf("writing secrets config file: %w", err)
	}

	return &api.AddVaultOutput{}, nil
}

func (ws *Workspace) DescribeVaults(ctx context.Context, input *api.DescribeVaultsInput) (*api.DescribeVaultsOutput, error) {
	vaultConfigs, err := ws.getVaultConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting vault configs: %w", err)
	}

	descriptions := make([]api.VaultDescription, len(vaultConfigs))
	for i, vaultConfig := range vaultConfigs {
		// TODO: add a status command.
		_, err = ws.EsvClient.GetWorkspaceSecrets(vaultConfig.url)
		descriptions[i] = api.VaultDescription{
			Name:      vaultConfig.name,
			URL:       vaultConfig.url,
			Connected: err == nil,
			NeedsAuth: errors.Is(err, esv.AuthError),
		}
	}

	return &api.DescribeVaultsOutput{Vaults: descriptions}, nil
}
