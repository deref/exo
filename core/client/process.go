// Generated file. DO NOT EDIT.

package client

import (
	"context"

	"github.com/deref/exo/core/api"
	josh "github.com/deref/exo/josh/client"
)

type Process struct {
	client *josh.Client
}

var _ api.Process = (*Process)(nil)

func NewProcess(client *josh.Client) *Process {
	return &Process{
		client: client,
	}
}

func (c *Process) Start(ctx context.Context, input *api.StartInput) (output *api.StartOutput, err error) {
	err = c.client.Invoke(ctx, "start", input, &output)
	return
}

func (c *Process) Stop(ctx context.Context, input *api.StopInput) (output *api.StopOutput, err error) {
	err = c.client.Invoke(ctx, "stop", input, &output)
	return
}
