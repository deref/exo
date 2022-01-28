// This file is for testing the event store in isolation.

package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/internal/eventd/api"
	"github.com/deref/exo/internal/eventd/sqlite"
	"github.com/deref/exo/internal/gensym"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/pier"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()

	// XXX These should not be necessary, the Josh machinery should not
	// be dependent on these, but rather should have some kind of middleware.
	ctx = logging.ContextWithLogger(ctx, logging.Default())
	ctx = telemetry.ContextWithTelemetry(ctx, &telemetry.Nop{})

	dbPath := "/tmp/exo-dev-eventd.sqlite3"
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		cmdutil.Fatalf("error opening sqlite db: %v", err)
	}
	defer db.Close()

	store := &sqlite.Store{
		DB:    db,
		IDGen: gensym.NewULIDGenerator(ctx),
	}

	// Commented out while transitioning to graphql implementation.
	//if err := store.Migrate(ctx); err != nil {
	//	cmdutil.Fatalf("error migrating: %v", err)
	//}

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
