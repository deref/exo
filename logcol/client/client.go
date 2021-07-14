// TODO: Generate this with JOSH tools.

package client

import (
	"context"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/logcol/api"
)

type LogCollector struct {
	client *josh.Client
}

var _ api.LogCollector = (*LogCollector)(nil)

func NewLogCollector(client *josh.Client) *LogCollector {
	return &LogCollector{client: client}
}

func (c *LogCollector) AddLog(ctx context.Context, input *api.AddLogInput) (output *api.AddLogOutput, err error) {
	err = c.client.Invoke(ctx, "add-log", input, &output)
	return
}

func (c *LogCollector) RemoveLog(ctx context.Context, input *api.RemoveLogInput) (output *api.RemoveLogOutput, err error) {
	err = c.client.Invoke(ctx, "remove-log", input, &output)
	return
}

func (c *LogCollector) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (output *api.DescribeLogsOutput, err error) {
	err = c.client.Invoke(ctx, "describe-logs", input, &output)
	return
}

func (c *LogCollector) GetEvents(ctx context.Context, input *api.GetEventsInput) (output *api.GetEventsOutput, err error) {
	err = c.client.Invoke(ctx, "get-events", input, &output)
	return
}
