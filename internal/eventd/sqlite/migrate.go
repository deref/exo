package sqlite

import (
	"context"
	"fmt"
)

func (sto *Store) Migrate(ctx context.Context) error {
	if _, err := sto.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS event (
			stream TEXT NOT NULL,
			id TEXT NOT NULL,
			timestamp INTEGER NOT NULL,
			message TEXT NOT NULL
		);`); err != nil {
		return fmt.Errorf("creating event table: %w", err)
	}
	if _, err := sto.DB.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS stream_event ON event ( stream, id )`); err != nil {
		return fmt.Errorf("creating stream_event index: %w", err)
	}
	if _, err := sto.DB.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS event_timestamp ON event ( timestamp )`); err != nil {
		return fmt.Errorf("creating event_timestamp index: %w", err)
	}
	return nil
}
