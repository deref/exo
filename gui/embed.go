// +build bundle

package gui

import (
	"context"
	"embed"
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
	return http.FileServer(http.FS(content))
}
