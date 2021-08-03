package main

import (
	"context"

	"github.com/deref/exo/config"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/logging"
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
	ctx := context.Background()

	ctx = config.WithConfig(ctx, cfg)

	logger := logging.Default()
	ctx = logging.ContextWithLogger(ctx, logger)

	return ctx
}

func main() {
	if err := config.LoadDefault(cfg); err != nil {
		cmdutil.Fatalf("loading config: %w", err)
	}
	knownPaths = cmdutil.MustMakeDirectories(cfg)
	if err := rootCmd.Execute(); err != nil {
		cmdutil.Fatalf("%w", err)
	}
}
