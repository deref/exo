package api

import (
	"context"

	"github.com/deref/graphql-go/encoding"
	"github.com/deref/graphql-go/gql"
)

type Service interface {
	Shutdown(context.Context) error
	// Execute a GraphQL operation synchronously.
	Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error
	// Schedule asynchronous execution of a GraphQL mutation.
	Enqueue(ctx context.Context, mutation string, vars map[string]interface{}) (jobID string, err error)
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
	})
	res := &encoding.SelectionUnmarshaler{
		Selection: sel,
	}
	return svc.Do(ctx, doc, vars, res)
}
