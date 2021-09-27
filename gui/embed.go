// +build bundle

package gui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/deref/exo/internal/config"
)

//go:embed dist/*
var dist embed.FS

func NewHandler(ctx context.Context, cfg config.GUIConfig) http.Handler {
	content, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return &tokenCookieMiddleware{
		URL:  fmt.Sprintf("http://localhost:%d/", cfg.Port),
		Next: http.FileServer(http.FS(content)),
	}
}
