package cli

import (
	"fmt"
	"strings"

	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/install"
	"github.com/deref/exo/internal/telemetry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade exo",
	Long:  `Upgrade exo to the latest version.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		tel := telemetry.FromContext(ctx)
		if !tel.IsEnabled() {
			fmt.Println("Cannot check current version - telemetry disabled.")
			if install.IsManaged {
				fmt.Println("Please upgrade using your system package manager.")
			} else {
				fmt.Println("You may upgrade with the following command.")
				fmt.Println("\tcurl -sL https://exo.deref.io/install | sh")
			}
			return nil
		}
		current := about.Version
		latest, err := tel.LatestVersion(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Current:\t%s\nLatest: \t%s\n", current, latest)

		switch strings.Compare(current, latest) {
		case 0:
			fmt.Println("You are already running the latest version")
		case -1:
			fmt.Println("Upgrade needed")
			inst := install.Get(cfg.VarDir + "/deviceid")
			// TODO: Prompt for confirmation?
			return inst.UpgradeSelf()
		case 1:
			fmt.Println("You are already running a prerelease version; not downgrading.")
		default:
			panic("Invalid version comparison")
		}
		return nil
	},
}
