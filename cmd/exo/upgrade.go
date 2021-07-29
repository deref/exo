package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/deref/exo/telemetry"
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
		ctx := newContext()
		if ok, reason := telemetry.CanSelfUpgrade(ctx); !ok {
			fmt.Fprintf(os.Stderr, "Self-upgrade disabled: %s\n", reason)
			os.Exit(1)
		}
		current := telemetry.CurrentVersion(ctx)
		latest, err := telemetry.LatestVersion(ctx)
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
			return telemetry.UpgradeSelf(ctx)
		case 1:
			fmt.Println("You are already running a prerelease version; not downgrading.")
		default:
			panic("Invalid version comparison")
		}
		return nil
	},
}
