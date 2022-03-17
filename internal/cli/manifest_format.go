package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestFormatCmd)
	addManifestFormatFlag(manifestFormatCmd, &manifestFormatFlags.Format)
}

var manifestFormatFlags struct {
	Format string
}

var manifestFormatCmd = &cobra.Command{
	Use:   "format [<path>]",
	Short: "Formats manifest file",
	Long: `Reformats a manifest file at the given path.

If path is not provided, uses the current workspace's configured manifest.

If the path is "-", reads from stdin and writes to stdout. Otherwise, the
specified file is updated in-place.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		vars := map[string]any{}
		if manifestFormatFlags.Format == "" {
			vars["format"] = (*string)(nil)
		} else {
			vars["format"] = manifestFormatFlags.Format
		}
		if len(args) > 0 && args[0] == "-" {
			content, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			var q struct {
				Manifest struct {
					Formatted string
				} `graphql:"makeManifest(content: $content, format: $format)"`
			}
			vars["content"] = string(content)
			if err := api.Query(ctx, svc, &q, vars); err != nil {
				return err
			}
			cmdutil.Show(q.Manifest.Formatted)
			return nil
		} else {
			vars["workspace"] = currentWorkspaceRef()
			if len(args) == 0 {
				vars["path"] = (*string)(nil)
			} else {
				vars["path"] = args[0]
			}
			var resp struct{}
			return svc.Do(ctx, &resp, `
				mutation ($workspace: String!, $format: String, $path: String) {
					formatManifest(workspace: $workspace, format: $format, path: $path) {
						__typename
					}
				}
			`, vars)
		}
	},
}
