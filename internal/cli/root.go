package cli

import (
	"context"

	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config = &config.Config{}
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

// newContext creates a global context that is used as the top-level process context.
// All request-specific contexts are derived from it.
func newContext() context.Context {
	ctx := context.Background()

	logger := logging.Default()
	ctx = logging.ContextWithLogger(ctx, logger)

	// This telemetry object will be replaced by one that includes the
	// device ID when the daemon starts. For one-off commands, it is
	// fine for the ID to not be populated.
	tel := telemetry.New(ctx, telemetry.Config{
		Disable: cfg.Telemetry.Disable,
	})
	ctx = telemetry.ContextWithTelemetry(ctx, tel)

	return ctx
}

func Main() {
	if err := config.LoadDefault(cfg); err != nil {
		cmdutil.Fatalf("loading config: %w", err)
	}
	if err := rootCmd.Execute(); err != nil {
		cmdutil.Fatalf("%w", err)
	}
}
