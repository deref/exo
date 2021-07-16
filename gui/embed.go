// +build bundle

package gui

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

func NewHandler(ctx context.Context) http.Handler {
	content, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(content))
}
