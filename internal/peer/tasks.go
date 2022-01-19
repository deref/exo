package peer

import (
	"context"
	"fmt"
	"os"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
)

func WorkTask(ctx context.Context, p *Peer, taskID string) error {
	workerId := fmt.Sprintf("peer:%d:", os.Getpid())

	var claim struct {
		Task struct {
			Mutation  string
			Variables string
		} `graphql:"claimTask(taskId: $taskId, workerId: $workerId)"`
	}
	if err := api.Mutate(ctx, p, &claim, map[string]interface{}{
		"taskId":   taskID,
		"workerId": workerId,
	}); err != nil {
		return fmt.Errorf("starting: %q", err)
	}
	task := claim.Task

	taskErr := func() error {
		var vars map[string]interface{}
		if err := jsonutil.UnmarshalString(task.Variables, &vars); err != nil {
			return fmt.Errorf("unmarshaling variables: %w", err)
		}

		mut := formatVoidMutation(task.Mutation, vars)

		tracker := &resolvers.TaskTracker{
			DB: p.db,
		}
		doCtx := resolvers.ContextWithTaskTracker(ctx, tracker)

		var res struct{}
		return p.Do(doCtx, mut, vars, &res)
	}()

	var finish struct {
		Void struct{} `graphql:"updateTask(id: $id, status: $status, finished: $finished, message: $message)"`
	}
	vars := map[string]interface{}{
		"id":       taskID,
		"finished": chrono.NowString(ctx),
	}
	if taskErr == nil {
		vars["status"] = api.TaskStatusSuccess
		vars["message"] = (*string)(nil)
	} else {
		vars["status"] = api.TaskStatusFailure
		vars["message"] = taskErr.Error()
	}
	if err := api.Mutate(ctx, p, &finish, vars); err != nil {
		return fmt.Errorf("finishing: %w", err)
	}
	return nil
}

func formatVoidMutation(mutation string, vars map[string]interface{}) string {
	arguments := make([]*ast.Argument, 0, len(vars))
	for k, v := range vars {
		arguments = append(arguments, &ast.Argument{
			Kind:  kinds.Argument,
			Name:  newNameNode(k),
			Value: newValueNode(v),
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
	return printer.Print(doc).(string)
}

func newNameNode(value string) *ast.Name {
	return &ast.Name{
		Kind:  kinds.Name,
		Value: value,
	}
}

func newValueNode(value interface{}) ast.Value {
	switch value := value.(type) {
	case string:
		return &ast.StringValue{
			Kind:  kinds.StringValue,
			Value: value,
		}
	default:
		panic(fmt.Errorf("cannot convert %T to GraphQL ast node", value))
	}
}
