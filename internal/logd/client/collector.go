// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/internal/josh/client"
	"github.com/deref/exo/internal/logd/api"
)

type LogCollector struct {
	client *josh.Client
}

var _ api.LogCollector = (*LogCollector)(nil)

func GetLogCollector(client *josh.Client) *LogCollector {
	return &LogCollector{
		client: client,
	}
}

func (c *LogCollector) ClearEvents(ctx context.Context, input *api.ClearEventsInput) (output *api.ClearEventsOutput, err error) {
	err = c.client.Invoke(ctx, "clear-events", input, &output)
	return
}

func (c *LogCollector) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (output *api.DescribeLogsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-logs", input, &output)
	return
}

func (c *LogCollector) AddEvent(ctx context.Context, input *api.AddEventInput) (output *api.AddEventOutput, err error) {
	err = c.client.Invoke(ctx, "add-event", input, &output)
	return
}

func (c *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = c.client.Invoke(ctx, "get-events", input, &output)
	return
}

func (c *LogCollector) RemoveOldEvents(ctx context.Context, input *api.RemoveOldEventsInput) (output *api.RemoveOldEventsOutput, err error) {
	err = c.client.Invoke(ctx, "remove-old-events", input, &output)
	return
}
