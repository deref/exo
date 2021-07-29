package telemetry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deref/exo"
	"github.com/deref/exo/config"
)

// CanSelfUpgrade checks whether this instance of exo is able to upgrade itself. If not,
// it returns a reason for the faulure.
// TODO: Find a better place for this and the autoinstalled/managed files to live.
func CanSelfUpgrade(ctx context.Context) (bool, string) {
	if isManaged {
		return false, "installed by a package manager"
	}
	if cfg, ok := config.GetConfig(ctx); ok && cfg.Telemetry.Disable {
		return false, "telemetry disabled - cannot check latest version"
	}

	return true, ""
}

// TrySelfUpgrade compares the installed and latest version, and if the installed
// version is behind, it replaces the exo binary with the latest version. If exo was
// installed by a package manager, this will not perform any version check or upgrade.
// The first return value indicates whether a newer version was downloaded.
func TrySelfUpgrade(ctx context.Context) (bool, error) {
	if ok, _ := CanSelfUpgrade(ctx); !ok {
		return false, nil
	}
	latest, err := LatestVersion(ctx)
	if err != nil {
		return false, err
	}
	if CurrentVersion(ctx) >= latest {
		return false, nil
	}

	if err := UpgradeSelf(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// LatestVersion returns the version of the running exo process.
func CurrentVersion(ctx context.Context) string {
	return exo.Version
}

// LatestVersion returns the latest version fetched from the web.
// TODO: Cache response!
func LatestVersion(ctx context.Context) (string, error) {
	resp, err := http.Get(exo.CheckVersionEndpoint)
	if err != nil {
		return "", fmt.Errorf("fetching latest version: %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading latest version: %w", err)
	}

	return string(body), nil
}
