package exod

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/deref/exo/gui"
	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/core/server"
	kernel "github.com/deref/exo/internal/core/server"
	"github.com/deref/exo/internal/core/state/statefile"
	"github.com/deref/exo/internal/esv"
	"github.com/deref/exo/internal/peer"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/task/api"
	taskserver "github.com/deref/exo/internal/task/server"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/token"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/exo/internal/util/sysutil"
	docker "github.com/docker/docker/client"
)

type Server struct {
	RedirectCrashDumps bool
}

// TODO: Some of the stuff in here should also be available to a worker daemon.
func (svr *Server) Run(ctx context.Context) {
	logger := logging.CurrentLogger(ctx)

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	MustMakeDirectories(cfg)

	if err := token.EnsureTokenFile(cfg.TokensFile); err != nil {
		cmdutil.Fatalf("ensuring token file: %w", err)
	}

	if svr.RedirectCrashDumps {
		// Panics will still write to stderr and some malbehaved code may write to
		// stdout or stderr. Redirect these file descriptors to truncated,
		// non-rotating, log files in the var directory. These logs won't be
		// preserved across runs, but can help us debug crashes where there is no
		// terminal attached.
		for _, redirect := range []struct {
			FD   int
			Name string
		}{
			{1, "stdout"},
			{2, "stderr"},
		} {
			dumpPath := filepath.Join(cfg.VarDir, "exod."+redirect.Name)
			dumpFile, err := os.OpenFile(dumpPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0600)
			if err != nil {
				logger.Infof("creating %s dump file: %v", redirect.Name, err)
			}
			if err := sysutil.Dup2(int(dumpFile.Fd()), redirect.FD); err != nil {
				logger.Infof("redirecting %s: %v", redirect.Name, err)
			}
		}
	}

	// When running as a daemon, we want to use the root filesystem to
	// avoid accidental relative path handling and to prevent tieing up
	// any mounted filesystem.
	if err := os.Chdir("/"); err != nil {
		cmdutil.Fatalf("chdir failed: %w", err)
	}

	statePath := filepath.Join(cfg.VarDir, "state.json")
	store := statefile.New(statePath)

	inst := about.GetInstall(filepath.Join(cfg.VarDir, "deviceid"))
	deviceID, err := inst.GetDeviceID()
	if err != nil {
		deviceID = "failed-to-set"
		logger.Infof("failed to initialize device: %v", err)
	}

	tel := telemetry.New(ctx, telemetry.Config{
		Disable:           cfg.Telemetry.Disable,
		DeviceID:          deviceID,
		DerefInternalUser: cfg.Telemetry.DerefInternalUser,
	})
	ctx = telemetry.ContextWithTelemetry(ctx, tel)
	tel.StartSession(ctx)
	tel.SendEvent(ctx, telemetry.SystemInfoIdentifiedEvent())

	service := &peer.Peer{
		SystemLog:   logger,
		VarDir:      cfg.VarDir,
		GUIEndpoint: fmt.Sprintf("http://localhost:%d", cfg.GUI.Port), // XXX should be constructed earlier than here.
		Debug:       true,                                             // XXX parameterize me.
	}
	if err := service.Init(ctx); err != nil {
		cmdutil.Fatalf("error initializing service: %v", err)
	}
	defer func() {
		if err := service.Shutdown(ctx); err != nil {
			logger.Infof("error shutting down service: %v", err)
		}
	}()

	dockerClient, err := docker.NewClientWithOpts()
	if err != nil {
		cmdutil.Fatalf("failed to create docker client: %v", err)
	}

	taskTracker := &task.TaskTracker{
		Store:  taskserver.NewTaskStore(),
		Logger: logger,
	}

	kernelCfg := &kernel.Config{
		Install:     inst,
		VarDir:      cfg.VarDir,
		Store:       store,
		SyslogPort:  cfg.Log.SyslogPort,
		Docker:      dockerClient,
		Logger:      logger,
		TaskTracker: taskTracker,
		TokenClient: cfg.GetTokenClient(),
		EsvClient:   esv.NewEsvClient(cfg.EsvTokenPath),
		ExoVersion:  about.Version,
		Service:     service,
	}

	// Commented out while transitioning to graphql implementation.
	//eventStore := &eventdsqlite.Store{
	//	DB:    db,
	//	IDGen: gensym.NewULIDGenerator(ctx),
	//}
	//
	//if err := eventStore.Migrate(ctx); err != nil {
	//	cmdutil.Fatalf("migrating event store: %v", err)
	//}
	//
	//syslogServer := &syslogd.Server{
	//	SyslogPort: kernelCfg.SyslogPort,
	//	Logger:     logger,
	//	Store:      eventStore,
	//}
	//ctx = log.ContextWithEventStore(ctx, eventStore)

	mux := server.BuildRootMux("/_exo/", kernelCfg)
	mux.Handle("/", gui.NewHandler(ctx, cfg.GUI))

	// TODO: Can these become cron components of a "system stack"?
	// Each execution would enqueue a job; somehow preserving backpressure.
	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()

		// Commented out while transitioning to graphql implementation.
		//go func() {
		//	for {
		//		select {
		//		case <-ctx.Done():
		//			return
		//		case <-time.After(5 * time.Second):
		//			if _, err := eventStore.RemoveOldEvents(ctx, &eventdapi.RemoveOldEventsInput{}); err != nil {
		//				logger.Infof("error removing old events: %v", err)
		//			}
		//		}
		//	}
		//}()

		// Commented out while transitioning to graphql implementation.
		//go func() {
		//	if err := syslogServer.Run(ctx); err != nil {
		//		cmdutil.Fatalf("syslog server error: %w", err)
		//	}
		//}()

		go func() {
			for {
				select {
				case <-ctx.Done():
				case <-time.After(10 * time.Second):
					if _, err := taskTracker.Store.EvictTasks(ctx, &api.EvictTasksInput{}); err != nil {
						logger.Infof("task eviction error: %w", err)
					}
				}
			}
		}()
	}

	validHosts := []string{fmt.Sprintf("localhost:%d", cfg.HTTPPort)}
	handler := httputil.HandlerWithContext(ctx, &httputil.HostAllowListHandler{
		Hosts: validHosts,
		Next:  mux,
	})

	addr := cmdutil.GetAddr(cfg)
	logger.Infof("listening for API calls at %s", addr)
	cmdutil.ListenAndServe(ctx, &http.Server{
		Addr:    addr,
		Handler: handler,
		// ErrorLog: XXX set this,
	})
}

func MustMakeDirectories(cfg *config.Config) {
	paths := []string{
		cfg.HomeDir,
		cfg.BinDir,
		cfg.VarDir,
		cfg.RunDir,
	}

	for _, path := range paths {
		if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
			cmdutil.Fatalf("making %q: %w", path, err)
		}
	}
}
