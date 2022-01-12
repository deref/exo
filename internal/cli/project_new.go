package cli

import (
	"fmt"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(projectNewCmd)
	projectNewCmd.Flags().StringVar(&projectNewFlags.DisplayName, "display-name", "", "Display name of project.")
}

var projectNewFlags struct {
	DisplayName string
}

var projectNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new project",
	Long:  `Creates a new project.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var m struct {
			Project struct {
				ID string
			} `graphql:"newProject(displayName: $displayName)"`
		}
		if err := cl.Mutate(ctx, &m, map[string]interface{}{
			"displayName": graphql.String(projectNewFlags.DisplayName),
		}); err != nil {
			return err
		}
		fmt.Println(m.Project.ID)
		return nil
	},
}