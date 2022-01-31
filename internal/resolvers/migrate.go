// XXX This is not graphql specific!

package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
	"github.com/mattn/go-sqlite3"
)

func (r *MutationResolver) Migrate(ctx context.Context) error {
	// Cluster.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cluster (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL
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
		INSERT INTO cluster ( id, name )
		VALUES ( ?, ? )
	`, gensym.RandomBase32(), "local",
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
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			spec TEXT NOT NULL,
			disposed TEXT
	);`); err != nil {
		return fmt.Errorf("creating component table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS
		component_stack_id_and_name ON component ( stack_id, name )
		WHERE disposed IS NULL
	`); err != nil {
		return fmt.Errorf("creating component_stack_id_and_name index: %w", err)
	}

	// Resource.

	// TODO: Consider dropping job_id field in favor of table for reified locks.
	// XXX remove message field; use events.
	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS resource (
			id TEXT NOT NULL PRIMARY KEY,
			type TEXT NOT NULL,
			iri TEXT,
			owner_type TEXT,
			owner_id TEXT,
			job_id TEXT,
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
			worker_id TEXT,
			status TEXT NOT NULL,
			error_message TEXT,
			created TEXT NOT NULL,
			updated TEXT NOT NULL,
			started TEXT,
			canceled TEXT,
			finished TEXT,
			progress_current INT,
			progress_total INT
	);`); err != nil {
		return fmt.Errorf("creating job table: %w", err)
	}

	// Stream.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stream (
			id TEXT PRIMARY KEY,
			source_type TEXT NOT NULL,
			source_id TEXT NOT NULL,
			created TEXT NOT NULL,
			truncated TEXT
		);
	`); err != nil {
		return fmt.Errorf("creating stream table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS stream_source
		ON stream ( source_type, source_id )
	`); err != nil {
		return fmt.Errorf("creating stream_source index: %w", err)
	}

	// Event.

	if _, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS event (
			ulid TEXT PRIMARY KEY,
			stream_id TEXT NOT NULL,
			message TEXT NOT NULL,
			tags TEXT NOT NULL
	);`); err != nil {
		return fmt.Errorf("creating event table: %w", err)
	}

	if _, err := r.DB.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS stream_event
		ON event ( stream_id, ulid )
	`); err != nil {
		return fmt.Errorf("creating stream_event index: %w", err)
	}

	return nil
}
