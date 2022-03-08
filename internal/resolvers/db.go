package resolvers

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func OpenDB(ctx context.Context, dbPath string) (*sqlx.DB, error) {
	connStr := dbPath

	// Fully serialize transactions. Hurts performance, but reasonable for
	// an embedded database, as long as transactions are kept small.
	// Helps dramatically with simplicity and correctness.
	connStr += "?_txlock=exclusive"

	return sqlx.Open("sqlite3", connStr)
}
