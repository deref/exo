package exod

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	golog "log"

	"github.com/deref/exo/gui"
	"github.com/deref/exo/internal/config"
	"github.com/deref/exo/internal/core/server"
	kernel "github.com/deref/exo/internal/core/server"
	"github.com/deref/exo/internal/core/state/statefile"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/logd"
	logdserver "github.com/deref/exo/internal/logd/server"
	"github.com/deref/exo/internal/logd/store/badger"
	"github.com/deref/exo/internal/providers/core/components/log"
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
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Main(ctx context.Context) {
	cmd, err := cmdutil.ParseArgs(os.Args)
	if err != nil {
		cmdutil.Fatalf("parsing arguments: %w", err)
	}

	RunServer(ctx, cmd.Flags)
}

func RunServer(ctx context.Context, flags map[string]string) {
	logger := logging.CurrentLogger(ctx)
	tel := telemetry.FromContext(ctx)

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	MustMakeDirectories(cfg)

	tokenClient, err := token.EnsureTokenFile(cfg.TokenFile)
	if err != nil {
		cmdutil.Fatalf("ensuring token file: %w", err)
	}

	tel.StartSession(ctx)

	_, forceStdLog := flags["force-std-log"]
	if !(forceStdLog || isatty.IsTerminal(os.Stdout.Fd())) {
		// Replace the standard logger with a logger writes to the var directory
		// and handles log rotation.
		golog.SetOutput(&lumberjack.Logger{
			Filename:   filepath.Join(cfg.VarDir, "exod.log"),
			MaxSize:    20, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})

		// Panics will still write to stderr and some malbehaved code may write to
		// stdout or stderr. Redirect these file descriptors to truncated,
		// non-rotating, log files in the var directory. These logs  won't be
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
	// and mounted filesystem.
	if err := os.Chdir("/"); err != nil {
		cmdutil.Fatalf("chdir failed: %w", err)
	}

	statePath := filepath.Join(cfg.VarDir, "state.json")
	store := statefile.New(statePath)

	dockerClient, err := docker.NewClientWithOpts()
	if err != nil {
		cmdutil.Fatalf("failed to create docker client: %v", err)
	}

	taskTracker := &task.TaskTracker{
		Store:  taskserver.NewTaskStore(),
		Logger: logger,
	}

	kernelCfg := &kernel.Config{
		VarDir:      cfg.VarDir,
		Store:       store,
		SyslogPort:  cfg.Log.SyslogPort,
		Docker:      dockerClient,
		Logger:      logger,
		TaskTracker: taskTracker,
	}

	logsDir := filepath.Join(cfg.VarDir, "logs")
	logStore, err := badger.Open(ctx, logger, logsDir)
	if err != nil {
		cmdutil.Fatalf("opening logs store: %w", err)
	}
	defer logStore.Close()

	logd := &logd.Service{
		SyslogPort: kernelCfg.SyslogPort,
		Logger:     logger,
		LogCollector: logdserver.LogCollector{
			IDGen: gensym.NewULIDGenerator(ctx),
			Store: logStore,
		},
	}
	ctx = log.ContextWithLogCollector(ctx, &logd.LogCollector)

	mux := server.BuildRootMux("/_exo/", kernelCfg, tokenClient)
	mux.Handle("/", gui.NewHandler(ctx, cfg.GUI))

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()

		go func() {
			if err := logd.Run(ctx); err != nil {
				cmdutil.Fatalf("log collector error: %w", err)
			}
		}()

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

	addr := cmdutil.GetAddr(cfg)
	logger.Infof("listening for API calls at %s", addr)
	cmdutil.ListenAndServe(ctx, &http.Server{
		Addr:    addr,
		Handler: httputil.HandlerWithContext(ctx, mux),
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
