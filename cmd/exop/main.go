package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/exod/server"
	"github.com/deref/exo/exod/state/statefile"
	"github.com/deref/exo/fifofum"
	josh "github.com/deref/exo/josh/client"
	logd "github.com/deref/exo/logd/client"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/httputil"
	"github.com/deref/pier"
)

func main() {
	// In development, exop starts a process by exec-ing a copy of itself with the "fifofum"
	// subcommand so that it always uses the latest veresion of the code rather than relying
	// on a prebuilt version of `fifofum` being on the path.
	if len(os.Args) > 1 && os.Args[1] == "fifofum" {
		fifofum.Main(fmt.Sprintf("%s %s", os.Args[0], os.Args[1]), os.Args[2:])
		return
	}

	ctx := context.Background()

	paths := cmdutil.MustMakeDirectories()

	statePath := filepath.Join(paths.VarDir, "state.json")
	store := statefile.New(statePath)

	cfg := &server.Config{
		VarDir: paths.VarDir,
		Store:  store,
	}

	ctx = log.ContextWithLogCollector(ctx, logd.GetLogCollector(&josh.Client{
		HTTP: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					dialer := net.Dialer{}
					sockPath := filepath.Join(cfg.VarDir, "logd.sock")
					return dialer.DialContext(ctx, "unix", sockPath)
				},
			},
		},
		URL: "http://unix",
	}))

	mux := server.BuildRootMux("/", cfg)
	pier.Main(httputil.HandlerWithContext(ctx, mux))
}
