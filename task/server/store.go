package server

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/task/api"
	"github.com/deref/exo/util/errutil"
)

type TaskStore struct {
	mx    sync.Mutex
	tasks map[string]*api.TaskDescription
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*api.TaskDescription),
	}
}

func (sto *TaskStore) DescribeTasks(ctx context.Context, input *api.DescribeTasksInput) (*api.DescribeTasksOutput, error) {
	sto.mx.Lock()
	defer sto.mx.Unlock()

	descriptions := make([]api.TaskDescription, 0, len(sto.tasks))
	for _, task := range sto.tasks {
		descriptions = append(descriptions, *task)
	}
	sort.Slice(descriptions, func(i, j int) bool {
		return strings.Compare(descriptions[i].ID, descriptions[j].ID) < 0
	})

	return &api.DescribeTasksOutput{
		Tasks: descriptions,
	}, nil
}

func (sto *TaskStore) CreateTask(ctx context.Context, input *api.CreateTaskInput) (output *api.CreateTaskOutput, err error) {
	sto.mx.Lock()
	defer sto.mx.Unlock()

	id := gensym.RandomBase32()
	desc := api.TaskDescription{
		ID:       id,
		ParentID: input.ParentID,
		Name:     input.Name,
		Status:   api.StatusPending,
		Created:  chrono.NowString(ctx),
		Updated:  chrono.NowString(ctx),
	}
	if desc.ParentID == nil {
		desc.JobID = id
	} else {
		parent, ok := sto.tasks[*desc.ParentID]
		if !ok {
			return nil, errutil.HTTPErrorf(http.StatusNotFound, "no such parent task: %q", *desc.ParentID)
		}
		desc.JobID = parent.JobID
	}

	return &api.CreateTaskOutput{
		ID: id,
	}, nil
}

func (sto *TaskStore) UpdateTask(ctx context.Context, input *api.UpdateTaskInput) (output *api.UpdateTaskOutput, err error) {
	sto.mx.Lock()
	defer sto.mx.Unlock()

	task, ok := sto.tasks[input.ID]
	if !ok {
		return nil, errutil.HTTPErrorf(http.StatusNotFound, "no such task: %q", input.ID)
	}

	if input.Status != nil {
		task.Status = *input.Status
	}
	if input.Message != nil {
		task.Message = *input.Message
	}
	if input.Started != nil {
		task.Started = input.Started
	}
	if input.Finished != nil {
		task.Finished = input.Finished
	}

	task.Updated = chrono.NowString(ctx)

	return &api.UpdateTaskOutput{}, nil
}

func (sto *TaskStore) EvictTasks(ctx context.Context, input *api.EvictTasksInput) (output *api.EvictTasksOutput, err error) {
	sto.mx.Lock()
	defer sto.mx.Unlock()

	expire := chrono.IsoNano(chrono.Now(ctx).Add(-1 * time.Hour))

	// Gather most recent update times for each job.
	jobs := make(map[string]string)
	for _, task := range sto.tasks {
		updated := jobs[task.JobID]
		if strings.Compare(task.Updated, updated) > 0 {
			jobs[task.JobID] = task.Updated
		}
	}

	// Remove any tasks for jobs last updated bevior the expiration time.
	for id, task := range sto.tasks {
		updated := jobs[task.JobID]
		if strings.Compare(updated, expire) < 0 {
			delete(jobs, id)
		}
	}

	return &api.EvictTasksOutput{}, nil
}
