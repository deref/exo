package exod

import (
	"context"
	"net/http"
	"os"
	"os/signal"

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

	ctx := kernel.NewContext(context.Background(), cfg)
	ctx, done := signal.NotifyContext(ctx, os.Interrupt)
	defer done()

	collector := logd.NewLogCollector(ctx, &logd.Config{
		VarDir: cfg.VarDir,
	})
	ctx = log.ContextWithLogCollector(ctx, collector)

	addr := cmdutil.GetAddr()

	mux := http.NewServeMux()

	kernelPattern := "/_exo/"
	mux.Handle(kernelPattern, kernel.NewHandler(ctx, &kernel.Config{
		VarDir:     cfg.VarDir,
		MuxPattern: kernelPattern,
	}))

	mux.Handle("/", gui.NewHandler(ctx))

	go func() {
		if err := collector.Run(ctx); err != nil {
			cmdutil.Warnf("log collector error: %w", err)
		}
	}()

	if err := http.ListenAndServe(addr, mux); err != nil {
		cmdutil.Fatalf("listening: %w", err)
	}
}
