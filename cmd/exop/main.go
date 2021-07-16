package main

import (
	"context"
	"net"
	"net/http"
	"path/filepath"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/components/log"
	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/client"
	"github.com/deref/pier"
)

func main() {
	paths := cmdutil.MustMakeDirectories()
	cfg := &server.Config{
		VarDir:     paths.VarDir,
		MuxPattern: "/",
	}
	ctx := server.NewContext(context.Background(), cfg)
	ctx = log.ContextWithLogCollector(ctx, logcol.NewLogCollector(&josh.Client{
		HTTP: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					dialer := net.Dialer{}
					sockPath := filepath.Join(cfg.VarDir, "logcol.sock")
					return dialer.DialContext(ctx, "unix", sockPath)
				},
			},
		},
		URL: "http://unix/",
	}))
	pier.Main(server.NewHandler(ctx, cfg))
}
