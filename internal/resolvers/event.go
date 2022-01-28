package resolvers

import (
	"context"
	"fmt"
)

type EventResolver struct {
	Q *RootResolver
	EventRow
}

type EventRow struct {
	ULID     ULID   `db:"ulid"`
	StreamID string `db:"stream_id"`
	Message  string `db:"message"`
	Tags     Tags   `db:"tags"`
}

func (r *EventResolver) ID() string {
	return r.ULID.String()
}

func (r *MutationResolver) CreateEvent(ctx context.Context, args struct {
	StreamID  string
	Timestamp *Instant
	Message   string
	Tags      *Tags
}) (*EventResolver, error) {
	row := EventRow{
		StreamID: args.StreamID,
		ULID:     r.mustNextULID(ctx),
		Message:  args.Message,
	}
	if args.Tags != nil {
		row.Tags = *args.Tags
	}
	if err := r.insertRow(ctx, "event", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &EventResolver{
		Q:        r,
		EventRow: row,
	}, nil
}
