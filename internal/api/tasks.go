package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/util/logging"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
	"golang.org/x/sync/errgroup"
)

type WorkerPool struct {
	Service      Service
	Concurrency  int
	WorkerPrefix string
	JobID        string
}

type Worker struct {
	Service Service
	ID      string
	// If non-empty, only work tasks from the specified job.  Otherwise, work any
	// available task.
	JobID string
	// Called when there is no more work expected.
	OnDone func()
}

func (pool *WorkerPool) Run(ctx context.Context) error {
	workCtx := ctx
	acquireCtx, stopAcquiring := context.WithCancel(workCtx)
	defer stopAcquiring()

	var eg errgroup.Group
	for i := 0; i < pool.Concurrency; i++ {
		workerID := fmt.Sprintf("%s:%d", pool.WorkerPrefix, i)
		eg.Go(func() error {
			worker := &Worker{
				Service: pool.Service,
				ID:      workerID,
				JobID:   pool.JobID,
				OnDone:  stopAcquiring,
			}
			return worker.Run(acquireCtx, workCtx)
		})
	}
	return eg.Wait()
}

// We need a separate acquire context, so that we can cancel task acquisition
// immediately without cancelling tasks that are currently being worked.
// TODO: Figure out some less awkward interface for this.
func (worker *Worker) Run(acquireCtx context.Context, workCtx context.Context) error {
	logger := logging.CurrentLogger(workCtx).Sublogger(fmt.Sprintf("worker %s", worker.ID))
	var jobID *string
	if worker.JobID != "" {
		jobID = &worker.JobID
	}
	acquireVars := map[string]interface{}{
		"workerId": worker.ID,
		"jobId":    jobID,
	}
	for {
		var acquired struct {
			Task *struct {
				ID       string
				Mutation string
			} `graphql:"acquireTask(workerId: $workerId, jobId: $jobId)"`
		}
		if err := Mutate(acquireCtx, worker.Service, &acquired, acquireVars); err != nil {
			return fmt.Errorf("acquiring task: %w", err)
		}
		task := acquired.Task
		if task == nil {
			if worker.OnDone != nil {
				worker.OnDone()
			}
			return nil
		}
		logger.Infof("acquired %s task: %s", task.Mutation, task.ID)
		// XXX The way this is currently set up, it doesn't make much sense for
		// acquireTask and startTask to be separate.
		err := worker.workTask(workCtx, task.ID)
		if err == nil {
			logger.Infof("completed task: %s", task.ID)
		} else {
			logger.Infof("task %s failure: %v", task.ID, err)
		}
	}
}

func (worker *Worker) workTask(ctx context.Context, id string) error {
	var start struct {
		Task struct {
			ID        string
			JobID     string
			Mutation  string
			Arguments JSONObject
		} `graphql:"startTask(id: $id, workerId: $workerId)"`
	}
	if err := Mutate(ctx, worker.Service, &start, map[string]interface{}{
		"id":       id,
		"workerId": worker.ID,
	}); err != nil {
		return fmt.Errorf("starting: %w", err)
	}
	task := start.Task

	var eg errgroup.Group

	// Poll for task cancellation and to act as worker heartbeat.
	waitCtx, stopPolling := context.WithCancel(ctx)
	defer stopPolling()
	eg.Go(func() error {
		defer stopPolling()
		vars := map[string]interface{}{
			"id":       id,
			"workerId": worker.ID,
		}
		for {
			var m struct {
				Task struct {
					ID       string
					Finished string
					Canceled string
				} `graphql:"updateTask(id: $id, workerId: $workerId)"`
			}
			if err := Mutate(ctx, worker.Service, &m, vars); err != nil {
				return fmt.Errorf("heartbeat failed: %w", err)
			}
			if m.Task.Canceled != "" || m.Task.Finished != "" {
				return nil
			}
			err := chrono.Sleep(waitCtx, time.Second)
			if errors.Is(err, context.Canceled) {
				return nil
			}
		}
	})

	// Do the actual work and record the results.
	var res struct{}
	eg.Go(func() error {
		defer stopPolling()

		taskErr := func() error {
			args := (map[string]interface{})(task.Arguments)
			mut, err := formatVoidMutation(task.Mutation, args)
			if err != nil {
				return fmt.Errorf("encoding mutation: %w", err)
			}

			taskCtx := ContextWithVariables(ctx, ContextVariables{
				TaskID:   task.ID,
				JobID:    task.JobID,
				WorkerID: worker.ID,
			})
			return worker.Service.Do(taskCtx, &res, mut, args)
		}()

		var finish struct {
			Void struct {
				Typename string `graphql:"__typename"`
			} `graphql:"finishTask(id: $id, error: $error)"`
		}
		vars := map[string]interface{}{
			"id": id,
		}
		if taskErr == nil {
			vars["error"] = (*string)(nil)
		} else {
			vars["error"] = taskErr.Error()
		}
		if err := Mutate(ctx, worker.Service, &finish, vars); err != nil {
			return fmt.Errorf("finishing: %w", err)
		}
		return nil
	})

	return eg.Wait()
}

func formatVoidMutation(mutation string, vars map[string]interface{}) (string, error) {
	arguments := make([]*ast.Argument, 0, len(vars))
	for k, v := range vars {
		value, err := newValueNode(v)
		if err != nil {
			return "", fmt.Errorf("encoding value of %q variable: %w", k, err)
		}
		arguments = append(arguments, &ast.Argument{
			Kind:  kinds.Argument,
			Name:  newNameNode(k),
			Value: value,
		})
	}
	doc := &ast.Document{
		Kind: kinds.Document,
		Definitions: []ast.Node{
			&ast.OperationDefinition{
				Kind:      kinds.OperationDefinition,
				Operation: "mutation",
				SelectionSet: &ast.SelectionSet{
					Kind: kinds.SelectionSet,
					Selections: []ast.Selection{
						&ast.Field{
							Kind: kinds.Field,
							Name: &ast.Name{
								Kind:  kinds.Name,
								Value: mutation,
							},
							Arguments: arguments,
							SelectionSet: &ast.SelectionSet{
								Kind: kinds.SelectionSet,
								Selections: []ast.Selection{
									&ast.Field{
										Kind: kinds.Field,
										Name: newNameNode("__typename"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return printer.Print(doc).(string), nil
}

func newNameNode(value string) *ast.Name {
	return &ast.Name{
		Kind:  kinds.Name,
		Value: value,
	}
}

func newValueNode(value interface{}) (ast.Value, error) {
	switch value := value.(type) {
	case string:
		return &ast.StringValue{
			Kind:  kinds.StringValue,
			Value: value,
		}, nil

	case int:
		return &ast.IntValue{
			Kind:  kinds.IntValue,
			Value: strconv.Itoa(value),
		}, nil

	case float64:
		return &ast.FloatValue{
			Kind:  kinds.FloatValue,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		}, nil

	default:
		return nil, fmt.Errorf("cannot convert %T to GraphQL ast node", value)
	}
}
