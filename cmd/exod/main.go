package main

import (
	"context"
	"net/http"

	"github.com/deref/exo"
)

func main() {
	ctx := exo.NewContext(context.Background())
	http.ListenAndServe(":3000", exo.NewHandler(ctx))
}
