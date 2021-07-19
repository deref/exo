package main

import "github.com/spf13/cobra"

var helpSubcmd = &cobra.Command{
	Use:   "help",
	Short: "Help about subcommand",
	Long:  `Help provides help for any subcommand in this command group.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Parent().Help()
	},
}
