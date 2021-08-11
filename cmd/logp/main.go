// For testing non-worker bits of logd.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/deref/exo/internal/config"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/logd"
	"github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/pier"
)

func main() {
	ctx := context.Background()
	ctx, done := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer done()

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	paths := cmdutil.MustMakeDirectories(cfg)

	ctx = telemetry.ContextWithTelemetry(ctx, telemetry.New(ctx, &config.TelemetryConfig{
		Disable: true,
	}))

	logd := &logd.Service{
		VarDir:     paths.VarDir,
		SyslogPort: cfg.Log.SyslogPort,
		Logger:     logging.Default(),
	}

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
