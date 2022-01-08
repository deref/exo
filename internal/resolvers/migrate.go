package resolvers

import (
	"context"
	"fmt"
)

func (r *MutationResolver) Migrate(ctx context.Context) error {
	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS project (
			id TEXT NOT NULL PRIMARY KEY,
			display_name TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating project table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS workspace (
			id TEXT NOT NULL PRIMARY KEY,
			project_id TEXT NOT NULL,
			root TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating workspace table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		workspace_root ON workspace ( root )
	`); err != nil {
		return fmt.Errorf("creating workspace_root index: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stack (
			id TEXT NOT NULL PRIMARY KEY,
			workspace_id TEXT
	);`); err != nil {
		return fmt.Errorf("creating stack table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		stack_workspace_id ON stack ( workspace_id )
	`); err != nil {
		return fmt.Errorf("creating stack_workspace_id index: %w", err)
	}

	return nil
}
