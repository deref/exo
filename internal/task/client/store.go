// Generated file. DO NOT EDIT.

package client

import (
	"context"

	josh "github.com/deref/exo/internal/josh/client"
	"github.com/deref/exo/internal/task/api"
)

type TaskStore struct {
	client *josh.Client
}

var _ api.TaskStore = (*TaskStore)(nil)

func GetTaskStore(client *josh.Client) *TaskStore {
	return &TaskStore{
		client: client,
	}
}

func (c *TaskStore) DescribeTasks(ctx context.Context, input *api.DescribeTasksInput) (output *api.DescribeTasksOutput, err error) {
	err = c.client.Invoke(ctx, "describe-tasks", input, &output)
	return
}

func (c *TaskStore) CreateTask(ctx context.Context, input *api.CreateTaskInput) (output *api.CreateTaskOutput, err error) {
	err = c.client.Invoke(ctx, "create-task", input, &output)
	return
}

func (c *TaskStore) UpdateTask(ctx context.Context, input *api.UpdateTaskInput) (output *api.UpdateTaskOutput, err error) {
	err = c.client.Invoke(ctx, "update-task", input, &output)
	return
}

func (c *TaskStore) EvictTasks(ctx context.Context, input *api.EvictTasksInput) (output *api.EvictTasksOutput, err error) {
	err = c.client.Invoke(ctx, "evict-tasks", input, &output)
	return
}
