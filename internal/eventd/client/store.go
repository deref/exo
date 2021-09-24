// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/eventd/api"
	josh "github.com/deref/exo/internal/josh/client"
)

type Store struct {
	client *josh.Client
}

var _ api.Store = (*Store)(nil)

func GetStore(client *josh.Client) *Store {
	return &Store{
		client: client,
	}
}

func (c *Store) ClearEvents(ctx context.Context, input *api.ClearEventsInput) (output *api.ClearEventsOutput, err error) {
	err = c.client.Invoke(ctx, "clear-events", input, &output)
	return
}

func (c *Store) DescribeStreams(ctx context.Context, input *api.DescribeStreamsInput) (output *api.DescribeStreamsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-streams", input, &output)
	return
}

func (c *Store) AddEvent(ctx context.Context, input *api.AddEventInput) (output *api.AddEventOutput, err error) {
	err = c.client.Invoke(ctx, "add-event", input, &output)
	return
}

func (c *Store) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = c.client.Invoke(ctx, "get-events", input, &output)
	return
}

func (c *Store) RemoveOldEvents(ctx context.Context, input *api.RemoveOldEventsInput) (output *api.RemoveOldEventsOutput, err error) {
	err = c.client.Invoke(ctx, "remove-old-events", input, &output)
	return
}
