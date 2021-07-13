package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/exod"
)

func main() {
	ctx := exod.NewContext(context.Background())
	http.ListenAndServe(":3000", exod.NewHandler(ctx))
}
