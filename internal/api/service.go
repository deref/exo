package api

import (
	"context"

	"github.com/deref/graphql-go/encoding"
	"github.com/deref/graphql-go/gql"
)

type Service interface {
	Shutdown(context.Context) error
	// Execute a GraphQL operation synchronously.
	// Implementations should also respect CurrentContextVariables.
	Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error
}

func Query(ctx context.Context, svc Service, q interface{}, vars map[string]interface{}) error {
	return doReflective(ctx, svc, gql.Query, q, vars)
}

func Mutate(ctx context.Context, svc Service, m interface{}, vars map[string]interface{}) error {
	return doReflective(ctx, svc, gql.Mutation, m, vars)
}

func doReflective(ctx context.Context, svc Service, typ gql.OperationType, sel interface{}, vars map[string]interface{}) error {
	doc := encoding.MustMarshalOperation(&gql.Operation{
		OperationDefinition: gql.OperationDefinition{
			Type:      typ,
			Selection: sel,
		},
		Variables: vars,
	})
	res := &encoding.SelectionUnmarshaler{
		Selection: sel,
	}
	return svc.Do(ctx, doc, vars, res)
}

// Schedule asynchronous execution of a GraphQL mutation.
func Enqueue(ctx context.Context, svc Service, mutation string, arguments map[string]interface{}) (jobID string, err error) {
	var m struct {
		Job struct {
			ID string
		} `graphql:"createTask(mutation: $mutation, arguments: $arguments)"`
	}
	err = Mutate(ctx, svc, &m, map[string]interface{}{
		"mutation":  mutation,
		"arguments": arguments,
	})
	return m.Job.ID, err
}
