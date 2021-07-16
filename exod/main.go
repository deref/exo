package exod

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/components/log"
	"github.com/deref/exo/gui"
	kernel "github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/server"
)

func Main() {
	paths := cmdutil.MustMakeDirectories()
	cfg := &kernel.Config{
		VarDir:     paths.VarDir,
		MuxPattern: "/",
	}
	ctx := kernel.NewContext(context.Background(), cfg)

	collector := logcol.NewLogCollector(ctx, &logcol.Config{
		VarDir: cfg.VarDir,
	})
	if err := collector.Start(ctx); err != nil {
		cmdutil.Fatalf("starting collector: %w", err)
	}
	ctx = log.ContextWithLogCollector(ctx, collector)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		collector.Stop(ctx)
		os.Exit(0)
	}()

	addr := cmdutil.GetAddr()

	mux := http.NewServeMux()

	kernelPattern := "/_exo/"
	mux.Handle(kernelPattern, kernel.NewHandler(ctx, &kernel.Config{
		VarDir:     cfg.VarDir,
		MuxPattern: kernelPattern,
	}))

	mux.Handle("/", gui.NewHandler(ctx))

	if err := http.ListenAndServe(addr, mux); err != nil {
		cmdutil.Fatalf("listening: %w", err)
	}
}
