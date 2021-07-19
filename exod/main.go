package exod

import (
	"context"
	"net/http"
	"os"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/gui"
	kernel "github.com/deref/exo/kernel/server"
	logd "github.com/deref/exo/logd/server"
	"github.com/deref/exo/util/cmdutil"
)

func Main() {
	paths := cmdutil.MustMakeDirectories()
	cfg := &kernel.Config{
		VarDir:     paths.VarDir,
		MuxPattern: "/",
	}

	if err := os.Chdir("/"); err != nil {
		cmdutil.Fatalf("chdir failed: %w", err)
	}

	ctx := kernel.NewContext(context.Background(), cfg)

	collector := logd.NewLogCollector(ctx, &logd.Config{
		VarDir: cfg.VarDir,
	})
	ctx = log.ContextWithLogCollector(ctx, collector)

	mux := http.NewServeMux()

	kernelPattern := "/_exo/"
	mux.Handle(kernelPattern, kernel.NewHandler(ctx, &kernel.Config{
		VarDir:     cfg.VarDir,
		MuxPattern: kernelPattern,
	}))

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
		Handler: mux,
	})
}
