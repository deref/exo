package cmdutil

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deref/exo/internal/util/logging"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, svr *http.Server) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(svr.ListenAndServe)
	eg.Go(func() error {
		ShutdownOnInterrupt(ctx, svr)
		return nil
	})
	if err := eg.Wait(); err != nil {
		Fatal(err)
	}
}

func Serve(ctx context.Context, l net.Listener, svr *http.Server) {
	go svr.Serve(l)
	ShutdownOnInterrupt(ctx, svr)
}

func ShutdownOnInterrupt(ctx context.Context, svr *http.Server) {
	logger := logging.CurrentLogger(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()

	timeout := 5 * time.Second
	logger.Infof("server shutting down (timeout: %v)", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		logger.Infof("server shutdown")
		cancel()
	}()
	_ = svr.Shutdown(ctx)
}
