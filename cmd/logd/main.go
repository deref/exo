// Separate logd service for testing in isolation. Unused for production
// deployments.  By default, binds a unix domain socket for easy discovery from
// exop.  If PORT environment variable is set, listens there instead for easy
// testing via curl or httpie.

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/logcol/api"
	"github.com/deref/exo/logcol/server"
)

func main() {
	ctx := context.Background()
	collector := server.NewLogCollector()
	collector.Start(ctx)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		collector.Stop(ctx)
		os.Exit(0)
	}()

	port := os.Getenv("PORT")
	var network, addr string
	if port == "" {
		network = "unix"
		addr = "./var/logcol.sock"
		_ = os.Remove(addr)
	} else {
		network = "tcp"
		addr = ":" + port
	}
	listener, err := net.Listen(network, addr)
	if err != nil {
		cmdutil.Fatalf("error listening: %v", err)
	}

	fmt.Println("listening at", addr)

	handler := api.NewLogCollectorMux("/", collector)
	server := http.Server{
		Handler: handler,
	}
	server.Serve(listener)
}
