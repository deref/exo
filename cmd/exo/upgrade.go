package main

import (
	"fmt"
	"strings"

	"github.com/deref/exo"
	"github.com/deref/exo/telemetry"
	"github.com/deref/exo/upgrade"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade exo",
	Long:  `Upgrade exo to the latest version.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		tel := telemetry.New(&cfg.Telemetry)
		current := exo.Version
		latest, err := tel.LatestVersion()
		if err != nil {
			return err
		}

		fmt.Printf("Current:\t%s\nLatest: \t%s\n", current, latest)

		switch strings.Compare(current, latest) {
		case 0:
			fmt.Println("You are already running the latest version")
		case -1:
			fmt.Println("Upgrade needed")
			// TODO: Prompt for confirmation?
			return upgrade.UpgradeSelf()
		case 1:
			fmt.Println("You are already running a prerelease version; not downgrading.")
		default:
			panic("Invalid version comparison")
		}
		return nil
	},
}
