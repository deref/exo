package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(jobCmd)
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Create and inspect jobs",
	Long:  `Contains subcommands for operating on jobs.`,
	Args:  cobra.NoArgs,
	// NOTE: Hidden because jobs are an internal feature currently. They are also
	// not yet scoped to workspaces. This may be promoted to a public feature
	// as the jobs system matures.
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
