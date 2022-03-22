package cli

import (
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	vaultCmd.AddCommand(vaultLSCmd)
	vaultLSCmd.Flags().BoolVarP(&vaultLSFlags.All, "all", "a", false, "Alias for --scope=all")
	vaultLSCmd.Flags().StringVar(&vaultLSFlags.Scope, "scope", "stack", "stack or all") // TODO: Support projects.
}

var vaultLSFlags struct {
	All   bool
	Scope string
}

var vaultLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List vaults",
	Long:  `Lists vaults in the given scope.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		scope := vaultLSFlags.Scope
		if vaultLSFlags.All {
			if !cmd.Flag("scope").Changed {
				scope = "all"
			} else if scope != "all" {
				return fmt.Errorf("--all conflicts with --scope=%q", scope)
			}
		}

		type vaultFragment struct {
			ID    string
			URL   string
			Error *string
		}
		var vaults []vaultFragment

		switch scope {
		case "all":
			var q struct {
				Vaults []vaultFragment `graphql:"allVaults"`
			}
			if err := api.Query(ctx, svc, &q, nil); err != nil {
				return err
			}
			vaults = q.Vaults

		case "stack":
			var q struct {
				Stack *struct {
					Vaults []vaultFragment
				} `graphql:"stackByRef(ref: $currentStack)"`
			}
			mustQueryStack(ctx, &q, nil)
			vaults = q.Stack.Vaults

		default:
			return fmt.Errorf("unknown scope: %q", resourceLSFlags.Scope)
		}

		w := cmdutil.NewTableWriter("ID", "URL", "STATUS")
		for _, vault := range vaults {
			status := "ok"
			if vault.Error != nil {
				status = *vault.Error
			}
			w.WriteRow(vault.ID, vault.URL, status)
		}
		w.Flush()
		return nil
	},
}
