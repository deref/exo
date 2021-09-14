// +build !managed

package upgrade

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/deref/exo"
)

const IsManaged = false

func UpgradeSelf(deviceID string) error {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return fmt.Errorf("creating temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	resp, err := http.Get(fmt.Sprintf("%s?id=%s&prev=%s", exo.UpdateScriptEndpoint, url.QueryEscape(deviceID), exo.Version))
	if err != nil {
		return fmt.Errorf("fetching update script: %w", err)
	}

	defer resp.Body.Close()
	if _, err := io.Copy(tmpfile, resp.Body); err != nil {
		return fmt.Errorf("saving update script: %w", err)
	}
	tmpfile.Close()

	if err := os.Chmod(tmpfile.Name(), 0700); err != nil {
		return fmt.Errorf("making script executable: %w", err)
	}

	cmd := exec.Command(tmpfile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running update script: %w", err)
	}

	return nil
}
