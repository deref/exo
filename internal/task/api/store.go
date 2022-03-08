// Generated file. DO NOT EDIT.

package api

import (
	"context"
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
	ID       string        `json:"id"`
	Status   *string       `json:"status"`
	Message  *string       `json:"message"`
	Started  *string       `json:"started"`
	Finished *string       `json:"finished"`
	Progress *TaskProgress `json:"progress"`
}

type UpdateTaskOutput struct {
}

type EvictTasksInput struct {
}

type EvictTasksOutput struct {
}

type TaskDescription struct {
	ID string `json:"id"`
	// ID of root task in this tree.
	JobID    string        `json:"jobId"`
	ParentID *string       `json:"parentId"`
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Created  string        `json:"created"`
	Updated  string        `json:"updated"`
	Started  *string       `json:"started"`
	Finished *string       `json:"finished"`
	Progress *TaskProgress `json:"progress"`
}

type TaskProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
