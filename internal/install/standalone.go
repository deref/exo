// +build !managed

package install

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/deref/exo/internal/about"
)

const IsManaged = false

func (i *Install) UpgradeSelf() error {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return fmt.Errorf("creating temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	deviceID, err := i.GetDeviceID()
	if err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf("%s?id=%s&prev=%s", about.UpdateScriptEndpoint, url.QueryEscape(deviceID), about.Version))
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
