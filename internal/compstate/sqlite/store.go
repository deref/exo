package sqlite

import (
	"context"
	"fmt"
	"math"

	"github.com/deref/exo/internal/compstate/api"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB *sqlx.DB
}

func (sto *Store) SetState(ctx context.Context, input *api.SetStateInput) (*api.SetStateOutput, error) {
	if input.ComponentID == "" {
		return nil, fmt.Errorf("invalid component id: %q", input.ComponentID)
	}
	tagsJSON := "{}"
	if input.Tags != nil {
		tagsJSON = jsonutil.MustMarshalString(input.Tags)
	}
	tx, err := sto.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()
	row := sto.DB.QueryRowContext(ctx, `
		INSERT INTO component_state (
			component_id, version,
			type, content, tags, timestamp
		)
		VALUES (
			?, COALESCE((
					SELECT MAX(version) + 1
					FROM component_state
					WHERE component_id = ?
				), 1),
			?, ?, ?, ?
		)
		RETURNING version;
	`, input.ComponentID, input.ComponentID, input.Type, input.Content, tagsJSON, input.Timestamp)
	var output api.SetStateOutput
	if err := row.Scan(&output.Version); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing: %w", err)
	}
	return &output, nil
}

func (sto *Store) GetStates(ctx context.Context, input *api.GetStatesInput) (*api.GetStatesOutput, error) {
	limit := mathutil.IntClamp(input.History, 1, 10)
	maxVersion := int(math.MaxInt32)
	if input.Version > 0 {
		maxVersion = input.Version
	}
	rows, err := sto.DB.QueryContext(ctx, `
		SELECT
			component_id,
			version,
			type,
			content,
			tags,
			timestamp
		FROM component_state
		WHERE component_id = ?
		AND version <= ?
		ORDER BY version DESC
		LIMIT ?
		`,
		input.ComponentID, maxVersion, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	output := api.GetStatesOutput{
		States: make([]api.State, 0, limit),
	}
	for rows.Next() {
		var state api.State
		var tags string
		if err := rows.Scan(&state.ComponentID, &state.Version, &state.Type, &state.Content, &tags, &state.Timestamp); err != nil {
			return nil, err
		}
		if err := jsonutil.UnmarshalString(tags, &state.Tags); err != nil {
			return nil, fmt.Errorf("unmarshalling state version %d tags: %w", state.Version, err)
		}
		output.States = append(output.States, state)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &output, nil
}
