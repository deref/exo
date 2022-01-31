package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type StreamResolver struct {
	Q *RootResolver
	StreamRow
}

type StreamRow struct {
	ID         string `db:"id"`
	SourceType string `db:"source_type"`
	SourceID   string `db:"source_id"`
}

func (r *QueryResolver) StreamByOwner(ctx context.Context, args struct {
	Type string
	ID   string
}) (*StreamResolver, error) {
	return r.streamByOwner(ctx, args.Type, args.ID)
}

func (r *QueryResolver) streamByOwner(ctx context.Context, typ string, id string) (*StreamResolver, error) {
	row := StreamRow{
		ID:         gensym.RandomBase32(),
		SourceType: typ,
		SourceID:   id,
	}
	if err := r.insertRowEx(ctx, "stream", &row, `
		ON CONFLICT (source_type, source_id)
		DO NOTHING
	`); err != nil {
		return nil, fmt.Errorf("upserting: %w", err)
	}
	return &StreamResolver{
		Q:         r,
		StreamRow: row,
	}, nil
}

func (r *StreamResolver) Message(ctx context.Context) (string, error) {
	return "", errors.New("TODO: StreamResolver.Message")
}
