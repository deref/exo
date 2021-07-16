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
	cfg := &server.Config{
		VarDir: cmdutil.MustVarDir(),
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

	if err := http.ListenAndServe("localhost:4000", server.NewHandler(ctx, cfg)); err != nil {
		cmdutil.Fatalf("listening: %w", err)
	}
}
