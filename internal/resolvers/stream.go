package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type StreamSource interface {
	Stream(context.Context) (*StreamResolver, error)
}

type StreamResolver struct {
	Q *RootResolver
	StreamRow
}

type StreamRow struct {
	ID         string `db:"id"`
	SourceType string `db:"source_type"`
	SourceID   string `db:"source_id"`
}

func (r *QueryResolver) streamById(ctx context.Context, id *string) (*StreamResolver, error) {
	stream := &StreamResolver{}
	err := r.getRowByKey(ctx, &stream.StreamRow, `
		SELECT *
		FROM stream
		WHERE id = ?
	`, id)
	if stream.ID == "" {
		stream = nil
	}
	return stream, err
}

func (r *MutationResolver) FindOrCreateStream(ctx context.Context, args struct {
	SourceType string
	SourceID   string
}) (*StreamResolver, error) {
	return r.findOrCreateStream(ctx, args.SourceType, args.SourceID)
}

func (r *MutationResolver) findOrCreateStream(ctx context.Context, sourceType string, sourceId string) (*StreamResolver, error) {
	row := StreamRow{
		ID:         gensym.RandomBase32(),
		SourceType: sourceType,
		SourceID:   sourceId,
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

func (r *StreamResolver) Source(ctx context.Context) (source StreamSource, err error) {
	var entity *Entity
	entity, err = r.Q.findEntity(ctx, r.SourceType, r.SourceID)
	if entity != nil {
		source = entity.Underlying.(StreamSource)
	}
	return
}

func (r *StreamResolver) Message(ctx context.Context) (string, error) {
	return "", errors.New("TODO: StreamResolver.Message")
}

func (r *StreamResolver) Events(ctx context.Context, args struct {
	Cursor *string
	Prev   *int32
	Next   *int32
	Filter *string
}) (*EventPageResolver, error) {
	q := eventQuery{
		StreamIDs: []string{r.ID},
	}
	if args.Cursor != nil {
		q.Cursor = *args.Cursor
	}
	if args.Prev != nil {
		q.Prev = int(*args.Prev)
	}
	if args.Next != nil {
		q.Next = int(*args.Next)
	}
	if args.Filter != nil {
		q.Filter = *args.Filter
	}
	return r.Q.findEvents(ctx, q)
}
