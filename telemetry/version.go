package telemetry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deref/exo/core"
)

// TODO: Find a better place for this and the autoinstalled/managed files to live.s
func CanSelfUpgrade() bool {
	return !isManaged
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
