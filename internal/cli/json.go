package cli

import (
	"encoding/json"
	"os"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().BoolVar(&jsonFlags.Indent, "indent", false, "output indented")
}

var jsonFlags struct {
	Indent bool
}

var jsonCmd = &cobra.Command{
	Use:    "json [properties...]",
	Hidden: true,
	Short:  "Builds a JSON object",
	Long: `Builds a JSON object of the properties specified as command line arguments.

Arguments pairs are expressed in one of two forms:

key=string
key:=raw

'key' and 'string' are unquoted JSON strings.
'raw' is an encoded JSON value.

Examples:

$ exo json name=Alice
{ "name": "Alice" }

$ exo json status:=404 'message:="not found"'
{ "status": 404, "message": "not found" }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		obj, err := cmdutil.ArgsToJsonObject(args)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		if jsonFlags.Indent {
			enc.SetIndent("", "  ")
		}
		return enc.Encode(obj)
	},
}
