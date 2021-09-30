package sqlite

import (
	"context"
	"fmt"
)

func (sto *Store) Migrate(ctx context.Context) error {
	if _, err := sto.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS component_state (
			component_id TEXT NOT NULL,
			version INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL,
			tags TEXT NOT NULL,
			timestamp INTEGER NOT NULL,

			PRIMARY KEY ( component_id, version )
		);`); err != nil {
		return fmt.Errorf("creating component_state table: %w", err)
	}
	return nil
}
