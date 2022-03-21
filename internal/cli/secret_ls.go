package cli

import (
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	secretCmd.AddCommand(secretLSCmd)
	secretLSCmd.Flags().BoolVar(&secretLSFlags.Reveal, "reveal", false, "Print secret values")
}

var secretLSFlags struct {
	Reveal bool
	// TODO: Explicit "scope" flag for consistency with other `ls` commands.
	Vault string
}

var secretLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "List secrets",
	Long: `Lists secrets in the given scope.
	
The default scope is all vaults attached to the current stack.

Secret values are not displayed unless --reveal is specified.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		type secretFragment struct {
			Name  string
			Value *string `graphql:"valueIf(reveal: $reveal)"`
			Vault struct {
				URL string
			}
		}
		var secrets []secretFragment
		columns := []string{"NAME"}
		if secretLSFlags.Reveal {
			columns = append(columns, "VALUE")
		}

		if secretLSFlags.Vault == "" {
			var q struct {
				Stack *struct {
					Secrets []secretFragment
				} `graphql:"stackByRef(ref: $currentStack)"`
			}
			mustQueryStack(ctx, &q, map[string]any{
				"reveal": secretLSFlags.Reveal,
			})
			secrets = q.Stack.Secrets
			columns = append(columns, "VAULT")
		} else {
			var q struct {
				Vault *struct {
					Secrets []secretFragment
				} `graphql:"vaultByRef(ref: $vault)"`
			}
			mustQueryStack(ctx, &q, map[string]any{
				"reveal": secretLSFlags.Reveal,
				"vault":  secretLSFlags.Vault,
			})
			secrets = q.Vault.Secrets
		}

		w := cmdutil.NewTableWriter(columns...)
		for _, secret := range secrets {
			data := map[string]string{
				"NAME":  secret.Name,
				"VAULT": secret.Vault.URL,
			}
			if secretLSFlags.Reveal {
				data["VALUE"] = *secret.Value
			}
			values := make([]string, len(columns))
			for i, column := range columns {
				values[i] = data[column]
			}
			w.WriteRow(values...)
		}
		w.Flush()
		return nil
	},
}
