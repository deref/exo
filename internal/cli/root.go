package cli

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/peer"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config = &config.Config{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&rootPersistentFlags.Cluster, "cluster", "", "")
}

var rootPersistentFlags struct {
	Cluster string
}

var rootCmd = &cobra.Command{
	Use:   "exo",
	Short: "Exo is a development environment process manager and log viewer.",
	Long: `A development environment process manager and log viewer.
For more information, see https://exo.deref.io`,
	// Automatic usage and error reporting behave badly, but Cobra Commander's
	// behavior is stable until v2.
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if !offline {
			// XXX client vs peer behavior.
			checkOrEnsureServer()
			var err error
			svc, err = peer.NewPeer(ctx, cfg.VarDir)
			if err != nil {
				return fmt.Errorf("initializing peer: %w", err)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if svc != nil {
			if err := svc.Shutdown(ctx); err != nil {
				cmdutil.Warnf("shutdown error: %w", err)
			}
			svc = nil
		}
		return nil
	},
}

// Most commands want to connect to the deamon, but those that don't, can set
// offline to true in their pre-run hook.  TODO: Implement a lazy client, so
// that offline doesn't need to be set explicitly.
var offline = false

// Will be initialized automatically, unless offline is true.
var svc api.Service

func Main() {
	ctx := context.Background()

	if err := config.LoadDefault(cfg); err != nil {
		cmdutil.Fatalf("loading config: %w", err)
	}

	logger := logging.Default()
	ctx = logging.ContextWithLogger(ctx, logger)

	// This telemetry object will be replaced by one that includes the
	// device ID when the daemon starts. For one-off commands, it is
	// fine for the ID to not be populated.
	tel := telemetry.New(ctx, telemetry.Config{
		Disable: cfg.Telemetry.Disable,
	})
	ctx = telemetry.ContextWithTelemetry(ctx, tel)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cmdutil.Fatalf("%w", err)
	}
}
