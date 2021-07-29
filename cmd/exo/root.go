package main

import (
	"context"
	"fmt"
	"os"

	"github.com/deref/exo/config"
	"github.com/deref/exo/util/cmdutil"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config = &config.Config{}
	knownPaths *cmdutil.KnownPaths
)

var rootCmd = &cobra.Command{
	Use:   "exo",
	Short: "Exo is a development environment process manager and log viewer.",
	Long: `A development environment process manager and log viewer.
For more information, see https://exo.deref.io`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func newContext() context.Context {
	return config.WithConfig(context.Background(), cfg)
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	knownPaths = cmdutil.MustMakeDirectories()
	if err := config.LoadDefault(cfg); err != nil {
		exitWithError(err)
	}
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}
