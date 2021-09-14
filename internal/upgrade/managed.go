// +build managed

package upgrade

import "fmt"

const IsManaged = true

func UpgradeSelf(deviceID string) error {
	fmt.Println(`This version of exo was installed via a package manager and does not support self-upgrade. Please use your package manager to update.

Alternatively, you may uninstall exo using your system's package manager and download the latest version with:

  curl -sL https://exo.deref.io/install | sh
`)
	return nil
}
