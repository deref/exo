// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

type TaskStore interface {
	DescribeTasks(context.Context, *DescribeTasksInput) (*DescribeTasksOutput, error)
	CreateTask(context.Context, *CreateTaskInput) (*CreateTaskOutput, error)
	UpdateTask(context.Context, *UpdateTaskInput) (*UpdateTaskOutput, error)
	EvictTasks(context.Context, *EvictTasksInput) (*EvictTasksOutput, error)
}

type DescribeTasksInput struct {

	// If supplied, filters tasks by job.
	JobIDs []string `json:"jobIds"`
}

type DescribeTasksOutput struct {
	Tasks []TaskDescription `json:"tasks"`
}

type CreateTaskInput struct {
	ParentID *string `json:"parentId"`
	Name     string  `json:"name"`
}

type CreateTaskOutput struct {
	ID    string `json:"id"`
	JobID string `json:"jobId"`
}

type UpdateTaskInput struct {
	ID       string  `json:"id"`
	Status   *string `json:"status"`
	Message  *string `json:"message"`
	Started  *string `json:"started"`
	Finished *string `json:"finished"`
}

type UpdateTaskOutput struct {
}

type EvictTasksInput struct {
}

type EvictTasksOutput struct {
}

func BuildTaskStoreMux(b *josh.MuxBuilder, factory func(req *http.Request) TaskStore) {
	b.AddMethod("describe-tasks", func(req *http.Request) interface{} {
		return factory(req).DescribeTasks
	})
	b.AddMethod("create-task", func(req *http.Request) interface{} {
		return factory(req).CreateTask
	})
	b.AddMethod("update-task", func(req *http.Request) interface{} {
		return factory(req).UpdateTask
	})
	b.AddMethod("evict-tasks", func(req *http.Request) interface{} {
		return factory(req).EvictTasks
	})
}

type TaskDescription struct {
	ID string `json:"id"`
	// ID of root task in this tree.
	JobID    string  `json:"jobId"`
	ParentID *string `json:"parentId"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Created  string  `json:"created"`
	Updated  string  `json:"updated"`
	Started  *string `json:"started"`
	Finished *string `json:"finished"`
}
