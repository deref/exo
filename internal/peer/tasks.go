package peer

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
)

func WorkTask(ctx context.Context, p *Peer, id string) error {
	workerId := fmt.Sprintf("peer:%d:", os.Getpid())

	var start struct {
		Task struct {
			Mutation  string
			Variables string
		} `graphql:"startTask(id: $id, workerId: $workerId)"`
	}
	if err := api.Mutate(ctx, p, &start, map[string]interface{}{
		"id":       id,
		"workerId": workerId,
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
