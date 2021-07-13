package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/core"
)

func main() {
	ctx := core.NewContext(context.Background())
	http.ListenAndServe(":3000", core.NewHandler(ctx))
}
