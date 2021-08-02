package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deref/exo/config"
	"github.com/deref/exo/core/server"
	"github.com/deref/exo/core/state/statefile"
	josh "github.com/deref/exo/josh/client"
	logd "github.com/deref/exo/logd/client"
	"github.com/deref/exo/providers/core/components/log"
	"github.com/deref/exo/supervise"
	"github.com/deref/exo/telemetry"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
	"github.com/deref/pier"
)

func main() {
	// In development, exop starts a process by exec-ing a copy of itself with the "supervise"
	// subcommand so that it always uses the latest veresion of the code rather than relying
	// on a prebuilt version of `supervise` being on the path.
	if len(os.Args) > 1 && os.Args[1] == "supervise" {
		// XXX: This is broken because supervise expects the syslod addr as the first argument.
		selfExec := os.Args[0]
		subCmd := os.Args[1]
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		subCmdArgs := os.Args[2:]
		supervise.Main(fmt.Sprintf("%s %s %s", selfExec, subCmd, wd), subCmdArgs)
		return
	}

	ctx := context.Background()

	cfg := &config.Config{}
	config.MustLoadDefault(cfg)
	paths := cmdutil.MustMakeDirectories(cfg)

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	serverCfg := &server.Config{
		VarDir:     paths.VarDir,
		Store:      store,
		Telemetry:  telemetry.New(&cfg.Telemetry),
		SyslogPort: log.SyslogPort,
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
