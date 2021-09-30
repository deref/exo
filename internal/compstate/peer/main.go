// This file is for testing the event store in isolation.

package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/internal/compstate/api"
	"github.com/deref/exo/internal/compstate/sqlite"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/pier"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()

	// SEE NOTE: [JOSH_CONTEXT].
	ctx = logging.ContextWithLogger(ctx, logging.Default())
	ctx = telemetry.ContextWithTelemetry(ctx, &telemetry.Nop{})

	dbPath := "/tmp/exo-dev-compstate.sqlite3"
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		cmdutil.Fatalf("error opening sqlite db: %v", err)
	}
	defer db.Close()

	store := &sqlite.Store{
		DB: db,
	}

	if err := store.Migrate(ctx); err != nil {
		cmdutil.Fatalf("error migrating: %v", err)
	}

	mb := josh.NewMuxBuilder("/")
	api.BuildStoreMux(mb, func(req *http.Request) api.Store {
		return store
	})
	mux := mb.Build()

	pier.Main(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(ctx)
		mux.ServeHTTP(w, req)
	}))
}
