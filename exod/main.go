package exod

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/exod/server"
	kernel "github.com/deref/exo/exod/server"
	"github.com/deref/exo/exod/state/statefile"
	"github.com/deref/exo/gui"
	logd "github.com/deref/exo/logd/server"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
)

func Main() {
	ctx := context.Background()

	paths := cmdutil.MustMakeDirectories()

	// When running as a deamon, we want to use the root filesystem to
	// avoid accidental relative path handling and to prevent tieing up
	// and mounted filesystem.
	if err := os.Chdir("/"); err != nil {
		cmdutil.Fatalf("chdir failed: %w", err)
	}

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	cfg := &kernel.Config{
		VarDir: paths.VarDir,
		Store:  store,
	}

	collector := logd.NewLogCollector(ctx, &logd.Config{
		VarDir: cfg.VarDir,
	})
	ctx = log.ContextWithLogCollector(ctx, collector)

	mux := server.BuildRootMux("/_exo/", cfg)
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
