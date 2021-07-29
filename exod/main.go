package exod

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	golog "log"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/config"
	"github.com/deref/exo/core/server"
	kernel "github.com/deref/exo/core/server"
	"github.com/deref/exo/core/state/statefile"
	"github.com/deref/exo/gui"
	"github.com/deref/exo/logd"
	"github.com/deref/exo/supervise"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
	"github.com/deref/exo/util/sysutil"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Main() {
	if len(os.Args) > 1 {
		subcommand := os.Args[1]
		switch subcommand {
		case "supervise":
			// XXX: This is broken because supervise expects the syslod addr as the first argument.
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			supervise.Main(fmt.Sprintf("%s %s %s", os.Args[0], subcommand, wd), os.Args[2:])
		case "server":
			RunServer()
		default:
			cmdutil.Fatalf("unknown subcommand: %q", subcommand)
		}
	} else {
		RunServer()
	}
}

func RunServer() {
	ctx := context.Background()

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	paths := cmdutil.MustMakeDirectories(cfg)

	if !isatty.IsTerminal(os.Stdout.Fd()) {
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
				golog.Printf("creating %s dump file: %v", redirect.Name, err)
			}
			if err := sysutil.Dup2(int(dumpFile.Fd()), redirect.FD); err != nil {
				golog.Printf("redirecting %s: %v", redirect.Name, err)
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

	kernelCfg := &kernel.Config{
		VarDir:     paths.VarDir,
		Store:      store,
		SyslogAddr: "localhost:4500", // XXX Configurable?
	}

	logd := &logd.Service{}
	logd.VarDir = kernelCfg.VarDir
	ctx = log.ContextWithLogCollector(ctx, &logd.LogCollector)

	mux := server.BuildRootMux("/_exo/", kernelCfg)
	mux.Handle("/", gui.NewHandler(ctx))

	{
		ctx, shutdown := context.WithCancel(ctx)
		defer shutdown()
		go func() {
			if err := logd.Run(ctx); err != nil {
				cmdutil.Fatalf("log collector error: %w", err)
			}
		}()
	}

	addr := cmdutil.GetAddr()
	golog.Printf("listening at %s", addr)

	cmdutil.ListenAndServe(ctx, &http.Server{
		Addr:    addr,
		Handler: httputil.HandlerWithContext(ctx, mux),
	})
}
