package api

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/graphql-go/encoding"
	"github.com/deref/graphql-go/gql"
)

type Service interface {
	Shutdown(context.Context) error
	// Execute a GraphQL operation synchronously.
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
func Enqueue(ctx context.Context, svc Service, mutation string, vars map[string]interface{}) (jobID string, err error) {
	var m struct {
		Job struct {
			ID string
		} `graphql:"createTask(mutation: $mutation, variables: $variables)"`
	}
	variablesJSON, err := jsonutil.MarshalString(vars)
	if err != nil {
		return "", fmt.Errorf("marshaling variables to json: %w", err)
	}
	err = Mutate(ctx, svc, &m, map[string]interface{}{
		"mutation":  mutation,
		"variables": variablesJSON,
	})
	return m.Job.ID, err
}
