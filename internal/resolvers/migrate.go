package resolvers

import (
	"context"
	"fmt"
)

func (r *MutationResolver) Migrate(ctx context.Context) error {
	// Project.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS project (
			id TEXT NOT NULL PRIMARY KEY,
			display_name TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating project table: %w", err)
	}

	// Workspace.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS workspace (
			id TEXT NOT NULL PRIMARY KEY,
			root TEXT NOT NULL,
			project_id TEXT
	);`); err != nil {
		return fmt.Errorf("creating workspace table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		workspace_root ON workspace ( root )
	`); err != nil {
		return fmt.Errorf("creating workspace_root index: %w", err)
	}

	// Stack.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stack (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			project_id TEXT,
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

	// Component.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS component (
			id TEXT NOT NULL PRIMARY KEY,
			stack_id TEXT NOT NULL,
			name TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating component table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS
		component_stack_id_and_name ON component ( stack_id, name )
	`); err != nil {
		return fmt.Errorf("creating component_stack_id_and_name index: %w", err)
	}

	// Resource.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS resource (
			iri TEXT NOT NULL PRIMARY KEY,
			component_id TEXT
	);`); err != nil {
		return fmt.Errorf("creating resource table: %w", err)
	}

	return nil
}
