package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/logcol/api"
	"github.com/deref/exo/logcol/server"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	lc := server.NewLogCollector()
	go lc.Collect(ctx, &api.CollectInput{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		os.Exit(0)
	}()
	http.ListenAndServe(":"+port, api.NewLogCollectorMux("/", lc))
}
