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
	"path/filepath"

	josh "github.com/deref/exo/josh/server"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/server"
	"github.com/deref/exo/util/cmdutil"
)

func main() {
	ctx := context.Background()
	ctx, done := signal.NotifyContext(ctx, os.Interrupt)
	defer done()

	paths := cmdutil.MustMakeDirectories()
	cfg := &server.Config{
		VarDir: paths.VarDir,
		Debug:  true,
	}

	collector := server.NewLogCollector(ctx, cfg)

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()
		go func() {
			if err := collector.Run(ctx); err != nil {
				cmdutil.Warnf("log collector error: %w", err)
			}
		}()
	}

	port := os.Getenv("PORT")
	var network, addr string
	if port == "" {
		network = "unix"
		addr = filepath.Join(cfg.VarDir, "logd.sock")
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

	muxb := josh.NewMuxBuilder("/")
	api.BuildLogCollectorMux(muxb, func(req *http.Request) api.LogCollector {
		return collector
	})
	cmdutil.Serve(ctx, listener, &http.Server{
		Handler: muxb.Build(),
	})
}
