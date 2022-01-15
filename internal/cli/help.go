package cli

import "github.com/spf13/cobra"

func makeHelpSubcmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Help about subcommand",
		Long:  `Help provides help for any subcommand in this command group.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			offline = true
			return cmd.Parent().PersistentPreRunE(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Parent().Help()
		},
	}
}
