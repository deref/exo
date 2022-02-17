package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(clusterEnvCmd)
}

var clusterEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment",
	Long: `Show cluster environment.

Finds a cluster as per the root "cluster" command, and prints exclusively that
cluster's environment, formats the environment with formatting as per the root
"env" command.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := lookupCluster(cmd)
		if err != nil {
			return err
		}
		showEnvironment(cluster.Environment)
		return nil
	},
}
