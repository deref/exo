package peer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
)

func RunWorker(ctx context.Context, p *Peer, workerID string, jobID *string) error {
	logger := logging.CurrentLogger(ctx)
	var timeout *int
	// Stop working tasks when a new task for this job hasn't appeared in a
	// while.  Note that this only works correclty when no other workers are
	// taking tasks for this job.
	// TODO: Would be better to acquire tasks until the job's status is finished.
	if jobID != nil {
		n := 1000
		timeout = &n
	}
	acquireVars := map[string]interface{}{
		"workerId": workerID,
		"jobId":    jobID,
		"timeout":  timeout,
	}
	for {
		var acquired struct {
			Task struct {
				ID string
			} `graphql:"acquireTask(workerId: $workerId, timeout: $timeout, jobId: $jobId)"`
		}
		if err := api.Mutate(ctx, p, &acquired, acquireVars); err != nil {
			return fmt.Errorf("acquiring task: %w", err)
		}
		taskID := acquired.Task.ID
		if taskID == "" {
			return nil
		}
		logger.Infof("acquired task: %s", taskID)
		err := WorkTask(ctx, p, workerID, taskID)
		if err == nil {
			logger.Infof("completed task: %s", taskID)
		} else {
			logger.Infof("task %s failure: %v", taskID, err)
		}
	}
}

func WorkTask(ctx context.Context, p *Peer, workerID string, id string) error {
	var start struct {
		Task struct {
			Mutation  string
			Variables string
		} `graphql:"startTask(id: $id, workerId: $workerId)"`
	}
	if err := api.Mutate(ctx, p, &start, map[string]interface{}{
		"id":       id,
		"workerId": workerID,
	}); err != nil {
		return fmt.Errorf("starting: %w", err)
	}
	task := start.Task

	taskErr := func() error {
		var vars map[string]interface{}
		if err := jsonutil.UnmarshalString(task.Variables, &vars); err != nil {
			return fmt.Errorf("unmarshaling variables: %w", err)
		}

		mut, err := formatVoidMutation(task.Mutation, vars)
		if err != nil {
			return fmt.Errorf("encoding mutation: %w", err)
		}

		doCtx := resolvers.ContextWithTaskID(ctx, id)
		var res struct{}
		return p.Do(doCtx, mut, vars, &res)
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
	if err := api.Mutate(ctx, p, &finish, vars); err != nil {
		return fmt.Errorf("finishing: %w", err)
	}
	return nil
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
