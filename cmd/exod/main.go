package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/kernel/server"
	"github.com/deref/exo/logcol/api"
	logcol "github.com/deref/exo/logcol/server"
)

func main() {
	ctx := server.NewContext(context.Background())

	lc := logcol.NewLogCollector()
	ctx = log.ContextWithLogCollector(ctx, lc)
	go lc.Collect(ctx, &api.CollectInput{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		os.Exit(0)
	}()

	http.ListenAndServe(":3000", server.NewHandler(ctx))
}
