package exod

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/components/process"
	"github.com/deref/exo/exod/server"
	kernel "github.com/deref/exo/exod/server"
	"github.com/deref/exo/exod/state/statefile"
	"github.com/deref/exo/gui"
	logd "github.com/deref/exo/logd/server"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
)

type Config struct {
	Fifofum process.FifofumConfig
}

func Main(cfg Config) {
	ctx := context.Background()

	paths := cmdutil.MustMakeDirectories()

	// When running as a daemon, we want to use the root filesystem to
	// avoid accidental relative path handling and to prevent tieing up
	// and mounted filesystem.
	if err := os.Chdir("/"); err != nil {
		cmdutil.Fatalf("chdir failed: %w", err)
	}

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	kernelCfg := &kernel.Config{
		VarDir:  paths.VarDir,
		Store:   store,
		Fifofum: cfg.Fifofum,
	}

	collector := logd.NewLogCollector(ctx, &logd.Config{
		VarDir: kernelCfg.VarDir,
	})
	ctx = log.ContextWithLogCollector(ctx, collector)

	mux := server.BuildRootMux("/_exo/", kernelCfg)
	mux.Handle("/", gui.NewHandler(ctx))

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()
		go func() {
			if err := collector.Run(ctx); err != nil {
				cmdutil.Warnf("log collector error: %w", err)
			}
		}()
	}

	cmdutil.ListenAndServe(ctx, &http.Server{
		Addr:    cmdutil.GetAddr(),
		Handler: httputil.HandlerWithContext(ctx, mux),
	})
}
