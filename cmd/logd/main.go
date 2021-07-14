// Separate logd service for testing in isolation. Unused for production deployments.

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
	collector := server.NewLogCollector()
	collector.Start(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		collector.Stop(ctx)
		os.Exit(0)
	}()
	http.ListenAndServe(":"+port, api.NewLogCollectorMux("/", collector))
}
