package exod

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	golog "log"

	"github.com/deref/exo/config"
	"github.com/deref/exo/core/server"
	kernel "github.com/deref/exo/core/server"
	"github.com/deref/exo/core/state/statefile"
	"github.com/deref/exo/gui"
	"github.com/deref/exo/logd"
	"github.com/deref/exo/providers/core/components/log"
	"github.com/deref/exo/supervise"
	"github.com/deref/exo/telemetry"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
	"github.com/deref/exo/util/logging"
	"github.com/deref/exo/util/sysutil"
	docker "github.com/docker/docker/client"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Main(ctx context.Context) {
	cmd, err := cmdutil.ParseArgs(os.Args)
	if err != nil {
		cmdutil.Fatalf("parsing arguments: %w", err)
	}

	subcommand := "server"
	if len(cmd.Args) > 0 {
		subcommand = cmd.Args[0]
		cmd.Args = cmd.Args[1:]
	}

	switch subcommand {
	case "supervise":
		// XXX: This is broken because supervise expects the syslod addr as the first argument.
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		supervise.Main(fmt.Sprintf("%s %s %s", cmd.Command, subcommand, wd), cmd.Args)

	case "server":
		RunServer(ctx, cmd.Flags)

	default:
		cmdutil.Fatalf("unknown subcommand: %q", subcommand)
	}
}

func RunServer(ctx context.Context, flags map[string]string) {
	logger := logging.CurrentLogger(ctx)

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	paths := cmdutil.MustMakeDirectories(cfg)

	tel := telemetry.New(&cfg.Telemetry)
	tel.StartSession()

	_, forceStdLog := flags["force-std-log"]
	if !(forceStdLog || isatty.IsTerminal(os.Stdout.Fd())) {
		// Replace the standard logger with a logger writes to the var directory
		// and handles log rotation.
		golog.SetOutput(&lumberjack.Logger{
			Filename:   filepath.Join(paths.VarDir, "exod.log"),
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
			dumpPath := filepath.Join(paths.VarDir, "exod."+redirect.Name)
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

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	dockerClient, err := docker.NewClientWithOpts()
	if err != nil {
		cmdutil.Fatalf("failed to create docker client: %v", err)
	}

	kernelCfg := &kernel.Config{
		VarDir:     paths.VarDir,
		Store:      store,
		Telemetry:  tel,
		SyslogPort: cfg.Log.SyslogPort,
		Docker:     dockerClient,
		Logger:     logger,
	}

	logd := &logd.Service{
		VarDir:     kernelCfg.VarDir,
		SyslogPort: kernelCfg.SyslogPort,
		Logger:     logger,
	}
	ctx = log.ContextWithLogCollector(ctx, &logd.LogCollector)

	mux := server.BuildRootMux("/_exo/", kernelCfg)
	mux.Handle("/", gui.NewHandler(ctx, cfg.GUI))

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()
		go func() {
			if err := logd.Run(ctx); err != nil {
				cmdutil.Fatalf("log collector error: %w", err)
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
