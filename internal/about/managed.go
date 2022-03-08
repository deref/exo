//go:build managed
// +build managed

package about

import "fmt"

const IsManaged = true

func (i *Install) UpgradeSelf() error {
	fmt.Println(`This version of exo was installed via a package manager and does not support self-upgrade. Please use your package manager to update.

Alternatively, you may uninstall exo using your system's package manager and download the latest version with:

  curl -sL https://exo.deref.io/install | sh
`)
	return nil
}
