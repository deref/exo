// XXX This is not graphql specific!

package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/scalars"
	"github.com/mattn/go-sqlite3"
)

func (r *MutationResolver) Migrate(ctx context.Context) error {
	now := scalars.Now(ctx)

	// Cluster.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cluster (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			environment_variables TEXT,
			updated TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating cluster table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		cluster_name ON cluster ( name )
	`); err != nil {
		return fmt.Errorf("creating cluster_name index: %w", err)
	}

	// TODO: Don't do this as part of migrate. SEE NOTE [DEFAULT_CLUSTER].
	if _, err := r.DB.ExecContext(ctx, `
		INSERT INTO cluster ( id, name, updated )
		VALUES ( ?, ?, ? )
	`, gensym.RandomBase32(), "local", now,
	); err != nil {
		var sqlErr sqlite3.Error
		if !(errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique) {
			return fmt.Errorf("inserting local cluster record: %w", err)
		}
	}

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
			project_id TEXT NOT NULL
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
			cluster_id TEXT NOT NULL,
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
			parent_id TEXT,
			type TEXT NOT NULL,
			name TEXT NOT NULL,
			key TEXT NOT NULL,
			spec TEXT NOT NULL,
			state TEXT NOT NULL,
			disposed TEXT
	);`); err != nil {
		return fmt.Errorf("creating component table: %w", err)
	}

	// TODO: Validate that these indexes are getting hit.

	// Nulls are distinct from themselves, so the stack serves as the parent id
	// for the purposes of indexing root components.
	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		component_path ON component ( stack_id, COALESCE(parent_id, stack_id), name )
		WHERE disposed IS NULL
	`); err != nil {
		return fmt.Errorf("creating component_path index: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		component_parent_id ON component ( parent_id )
		WHERE disposed IS NULL
	`); err != nil {
		return fmt.Errorf("creating component_parent_id index: %w", err)
	}

	// Resource.

	// TODO: Consider dropping task_id field in favor of table for reified locks.
	// XXX remove message field; use events.
	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS resource (
			id TEXT NOT NULL PRIMARY KEY,
			type TEXT NOT NULL,
			iri TEXT,
			owner_type TEXT,
			owner_id TEXT,
			task_id TEXT,
			model TEXT NOT NULL,
			status INT,
			message TEXT
	);`); err != nil {
		return fmt.Errorf("creating resource table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS
		resource_iri ON resource ( iri )
	`); err != nil {
		return fmt.Errorf("creating cluster_name index: %w", err)
	}

	// Task.
	// TODO: Tasks should be associated with a workspace, etc.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS task (
			id TEXT NOT NULL PRIMARY KEY,
			job_id TEXT NOT NULL,
			parent_id TEXT,
			mutation TEXT NOT NULL,
			arguments TEXT NOT NULL,
			key TEXT,
			worker_id TEXT,
			created TEXT NOT NULL,
			updated TEXT NOT NULL,
			started TEXT,
			canceled TEXT,
			finished TEXT,
			completed TEXT,
			progress_current INT NOT NULL,
			progress_total INT NOT NULL,
			error TEXT
	);`); err != nil {
		return fmt.Errorf("creating job table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		task_key ON task ( mutation, key )
		WHERE key IS NOT NULL
		AND completed IS NULL
	`); err != nil {
		return fmt.Errorf("creating task_parent_id index: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS
		task_parent_id ON task ( parent_id )
	`); err != nil {
		return fmt.Errorf("creating task_parent_id index: %w", err)
	}

	// Event.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS event (
			ulid TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			tags TEXT NOT NULL,
			source_type TEXT NOT NULL,
			workspace_id TEXT,
			stack_id TEXT,
			component_id TEXT,
			job_id TEXT,
			task_id TEXT
	);`); err != nil {
		return fmt.Errorf("creating event table: %w", err)
	}

	for _, related := range []string{
		"workspace",
		"stack",
		"component",
		"job",
		"task",
	} {
		// TODO: Validate that these indexes are getting hit.
		if _, err := r.DB.ExecContext(ctx, fmt.Sprintf(`
			CREATE UNIQUE INDEX IF NOT EXISTS %s_event
			ON event ( %s_id, ulid )
			WHERE %s_id IS NOT NULL
		`, related, related, related)); err != nil {
			return fmt.Errorf("creating %s_event index: %w", related, err)
		}
	}

	if _, err := r.DB.ExecContext(ctx, fmt.Sprintf(`
		CREATE UNIQUE INDEX IF NOT EXISTS system_event
		ON event ( ulid )
		WHERE source_type = 'System'
	`)); err != nil {
		return fmt.Errorf("creating system_event index: %w", err)
	}

	return nil
}
