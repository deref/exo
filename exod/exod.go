package exod

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/components/log"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/server"
)

func Main() {
	paths := cmdutil.MustMakeDirectories()
	cfg := &server.Config{
		VarDir: paths.VarDir,
	}
	ctx := server.NewContext(context.Background(), cfg)

	collector := logcol.NewLogCollector(ctx, &logcol.Config{
		VarDir: cfg.VarDir,
	})
	collector.Start(ctx)
	ctx = log.ContextWithLogCollector(ctx, collector)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		collector.Stop(ctx)
		os.Exit(0)
	}()

	addr := cmdutil.GetAddr()

	if err := http.ListenAndServe(addr, server.NewHandler(ctx, cfg)); err != nil {
		cmdutil.Fatalf("listening: %w", err)
	}
}
