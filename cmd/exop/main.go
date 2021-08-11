package main

import (
	"context"
	"net"
	"net/http"
	"path/filepath"

	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/core/server"
	"github.com/deref/exo/internal/core/state/statefile"
	josh "github.com/deref/exo/internal/josh/client"
	logd "github.com/deref/exo/internal/logd/client"
	"github.com/deref/exo/internal/providers/core/components/log"
	"github.com/deref/exo/internal/task"
	taskserver "github.com/deref/exo/internal/task/server"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/pier"
	docker "github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	paths := cmdutil.MustMakeDirectories(cfg)

	logger := logging.Default()
	ctx = logging.ContextWithLogger(ctx, logger)

	tel := telemetry.New(ctx, &cfg.Telemetry)
	ctx = telemetry.ContextWithTelemetry(ctx, tel)

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	dockerClient, err := docker.NewClientWithOpts()
	if err != nil {
		cmdutil.Fatalf("failed to create docker client: %v", err)
	}

	taskTracker := &task.TaskTracker{
		Store:  taskserver.NewTaskStore(),
		Logger: logger,
	}

	serverCfg := &server.Config{
		VarDir:      paths.VarDir,
		Store:       store,
		Logger:      logger,
		SyslogPort:  cfg.Log.SyslogPort,
		Docker:      dockerClient,
		TaskTracker: taskTracker,
	}

	ctx = log.ContextWithLogCollector(ctx, logd.GetLogCollector(&josh.Client{
		HTTP: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					dialer := net.Dialer{}
					sockPath := filepath.Join(serverCfg.VarDir, "logd.sock")
					return dialer.DialContext(ctx, "unix", sockPath)
				},
			},
		},
		URL: "http://unix",
	}))

	mux := server.BuildRootMux("/", serverCfg)
	pier.Main(httputil.HandlerWithContext(ctx, mux))
}
