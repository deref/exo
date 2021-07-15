package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/server"
)

func main() {
	cfg := &server.Config{
		VarDir: "./var", // XXX
	}
	ctx := server.NewContext(context.Background(), cfg)

	collector := logcol.NewLogCollector(&logcol.Config{
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

	http.ListenAndServe(":3000", server.NewHandler(ctx, cfg))
}
