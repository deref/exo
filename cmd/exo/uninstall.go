package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Hidden: true, // TODO: Is there a "less hidden" mode for uncommon commands?
	Use:    "uninstall",
	Short:  "Uninstall exo",
	Long: `Shutsdown everything, the uninstalls exo.

Performs the following steps:
- Stops and destroys all workspaces.
- Shutsdown the exo daemon.
- Deletes the exo home directory.

If this command fails, see <https://github.com/deref/exo/tree/main/doc/uninstall.md>
for manual uninstall instructions.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()

		cl := newClient()

		// Destroy workspaces.
		workspaces, err := cl.Kernel().DescribeWorkspaces(ctx, &api.DescribeWorkspacesInput{})
		if err != nil {
			cmdutil.Fatalf("describing workspaces: %w", err)
		}
		for _, workspace := range workspaces.Workspaces {
			fmt.Printf("destroying workspace: %s\n", workspace.ID)
			if _, err := cl.GetWorkspace(workspace.ID).Destroy(ctx, &api.DestroyInput{}); err != nil {
				return fmt.Errorf("destroying workspace %q: %w", workspace.ID, err)
			}
		}

		// Exit daemon.
		if err := killExod(); err != nil {
			fmt.Println("exiting daemon")
			return fmt.Errorf("exiting daemon: %w", err)
		}

		// Remove home directory.
		fmt.Printf("removing home directory: %s\n", knownPaths.ExoDir)
		if err := os.RemoveAll(knownPaths.ExoDir); err != nil {
			return fmt.Errorf("removing exo directory: %w", err)
		}

		fmt.Println("uninstalled!")
		return nil
	},
}
