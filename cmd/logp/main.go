// For testing non-worker bits of logd.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	josh "github.com/deref/exo/josh/server"
	"github.com/deref/exo/logd"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/pier"
)

func main() {
	ctx := context.Background()
	ctx, done := signal.NotifyContext(ctx, os.Interrupt)
	defer done()

	paths := cmdutil.MustMakeDirectories()

	logd := &logd.Service{}
	logd.VarDir = paths.VarDir

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()
		go func() {
			if err := logd.Run(ctx); err != nil {
				cmdutil.Warnf("log collector error: %w", err)
			}
		}()
	}

	muxb := josh.NewMuxBuilder("/")
	api.BuildLogCollectorMux(muxb, func(req *http.Request) api.LogCollector {
		return &logd.LogCollector
	})
	pier.Main(muxb.Build())
}
