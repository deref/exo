package telemetry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deref/exo/core"
)

// TODO: Find a better place for this and the autoinstalled/managed files to live.
func CanSelfUpgrade() bool {
	return !isManaged
}

// TrySelfUpgrade compares the installed and latest version, and if the installed
// version is behind, it replaces the exo binary with the latest version. If exo was
// installed by a package manager, this will not perform any version check or upgrade.
// The first return value indicates whether a newer version was downloaded.
func TrySelfUpgrade() (bool, error) {
	if !CanSelfUpgrade() {
		return false, nil
	}
	latest, err := LatestVersion()
	if err != nil {
		return false, err
	}
	// if CurrentVersion() >= latest {
	// 	return false, nil
	// }
	_ = latest

	if err := UpgradeSelf(); err != nil {
		return false, err
	}

	return true, nil
}

// LatestVersion returns the version of the running exo process.
func CurrentVersion() string {
	return core.Version
}

// LatestVersion returns the latest version fetched from the web.
// TODO: Cache response!
func LatestVersion() (string, error) {
	resp, err := http.Get(core.CheckVersionEndpoint)
	if err != nil {
		return "", fmt.Errorf("fetching latest version: %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading latest version: %w", err)
	}

	return string(body), nil
}
