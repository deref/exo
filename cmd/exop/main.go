package main

import (
	"context"
	"net"
	"net/http"
	"path/filepath"

	"github.com/deref/exo/components/log"
	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/client"
	"github.com/deref/pier"
)

func main() {
	ctx := server.NewContext(context.Background())
	varDir := "./var" // XXX parameterize and pass this around.
	ctx = log.ContextWithLogCollector(ctx, logcol.NewLogCollector(&josh.Client{
		HTTP: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
					dialer := net.Dialer{}
					return dialer.DialContext(ctx, "unix", filepath.Join(varDir, "logcol.sock"))
				},
			},
		},
		URL: "http://unix/",
	}))
	pier.Main(server.NewHandler(ctx))
}
