package cli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a workspace",
	Long: `Creates a workspace in the current directory and initializes it.

If a manifest file of a supported file is found, it will be applied. Otherwise,
a exo.hcl manifest file will be created.

Prints instructions for using the newly created workspace.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl := newClient()
		kernel := cl.Kernel()

		// Create workspace.
		root := cmdutil.MustGetwd()
		createOutput, err := cl.Kernel().CreateWorkspace(ctx, &api.CreateWorkspaceInput{
			Root: root,
		})
		if err != nil {
			return fmt.Errorf("creating workspace: %w", err)
		}
		workspaceID := createOutput.ID
		workspace := cl.GetWorkspace(workspaceID)
		fmt.Println("Workspace created.")

		// Search for existing manifest.
		fmt.Println()
		manifestOutput, err := workspace.ResolveManifest(ctx, &api.ResolveManifestInput{})
		if err != nil {
			return fmt.Errorf("resolving manifest: %w", err)
		}

		if manifestOutput.Path == "" {
			// Create a new empty manifest.
			manifestPath := filepath.Join(root, "exo.hcl")
			if err := ioutil.WriteFile(manifestPath, []byte(exohcl.Starter), 0600); err != nil {
				return fmt.Errorf("writing manifest file: %w", err)
			}
			fmt.Printf("Wrote manifest: %q\n", manifestPath)

			// Print instructions assuming there are no processes.
			fmt.Println(`
Ready to go! Here are some potential next steps:

Launch the graphical UI in a web browser:

		exo gui

Explore the CLI interface:

		exo help`)
		} else {
			// Apply an existing manifest.
			fmt.Println("Applying existing manifest...")
			if err := apply(ctx, kernel, workspace, []string{}); err != nil {
				return fmt.Errorf("applying manifest: %w", err)
			}

			// Print instructions assuming the manifest is non-empty.
			fmt.Println(`
Ready to go! Here are some potential next steps:

Launch graphical UI in web browser:

		exo gui

View logs in terminal:

		exo logs

Stop all processes in this workspace:

		exo stop`)
		}
		return nil
	},
}
