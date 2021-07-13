package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/kernel"
)

func main() {
	ctx := kernel.NewContext(context.Background())
	http.ListenAndServe(":3000", kernel.NewHandler(ctx))
}
