package cli

import (
	"context"
	"sync"

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
	rootCmd.PersistentFlags().BoolVar(&rootPersistentFlags.Async, "async", false, "Do not await long-running tasks")
}

var rootPersistentFlags struct {
	Cluster string
	Async   bool
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

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

	// XXX TODO: if client, do checkOrEnsureServer() instead of peer.NewPeer
	svc = newLazyService(func() api.Service {
		svc, err := peer.NewPeer(ctx, peer.PeerConfig{
			VarDir:      cfg.VarDir,
			GUIEndpoint: effectiveServerURL(),
		})
		if err != nil {
			cmdutil.Fatalf("initializing peer: %w", err)
		}
		return svc
	})
	defer func() {
		if err := svc.Shutdown(ctx); err != nil {
			cmdutil.Warnf("shutdown error: %w", err)
		}
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cmdutil.Fatalf("error: %w", err)
	}
}

type lazyService struct {
	mx    sync.Mutex
	thunk func() api.Service
	value api.Service
}

func newLazyService(thunk func() api.Service) *lazyService {
	return &lazyService{
		thunk: thunk,
	}
}

func (svc *lazyService) force() {
	svc.mx.Lock()
	defer svc.mx.Unlock()
	if svc.thunk == nil {
		return
	}
	thunk := svc.thunk
	svc.thunk = nil
	svc.value = thunk()
}

func (svc *lazyService) Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error {
	svc.force()
	return svc.value.Do(ctx, doc, vars, res)
}

func (svc *lazyService) Shutdown(ctx context.Context) error {
	svc.mx.Lock()
	defer svc.mx.Unlock()

	if svc.value == nil {
		return nil
	}
	err := svc.value.Shutdown(ctx)
	return err
}
