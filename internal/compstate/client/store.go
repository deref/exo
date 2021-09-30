// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/internal/compstate/api"
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

func (c *Store) SetState(ctx context.Context, input *api.SetStateInput) (output *api.SetStateOutput, err error) {
	err = c.client.Invoke(ctx, "set-state", input, &output)
	return
}

func (c *Store) GetStates(ctx context.Context, input *api.GetStatesInput) (output *api.GetStatesOutput, err error) {
	err = c.client.Invoke(ctx, "get-states", input, &output)
	return
}
