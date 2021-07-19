// +build !managed

package telemetry

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/deref/exo/core"
)

func UpgradeSelf() error {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		fmt.Errorf("creating temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	resp, err := http.Get(core.UpdateScriptEndpoint)
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

	out, err := exec.Command(tmpfile.Name()).Output()
	if err != nil {
		return fmt.Errorf("running update script: %w", err)
	}
	fmt.Printf(string(out))

	return nil
}
