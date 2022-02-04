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
	Do(ctx context.Context, res interface{}, doc string, vars map[string]interface{}) error
	// Begin a GraphQL subscription operation. Replaces res each time
	// subscription.Next returns. Consumer must call Stop() to cleanup.
	Subscribe(ctx context.Context, res interface{}, doc string, vars map[string]interface{}) Subscription
}

type Subscription interface {
	// Yields true each time an event is received and Service.Subscribe's res
	// argument is replaced. Closes when the subscription ends.
	C() <-chan bool
	Err() error
	Stop()
}

func Query(ctx context.Context, svc Service, q interface{}, vars map[string]interface{}) error {
	res, doc := buildReflectiveOperation(gql.Query, q, vars)
	return svc.Do(ctx, res, doc, vars)
}

func Mutate(ctx context.Context, svc Service, m interface{}, vars map[string]interface{}) error {
	res, doc := buildReflectiveOperation(gql.Mutation, m, vars)
	return svc.Do(ctx, res, doc, vars)
}

func Subscribe(ctx context.Context, svc Service, s interface{}, vars map[string]interface{}) Subscription {
	res, doc := buildReflectiveOperation(gql.Subscription, s, vars)
	return svc.Subscribe(ctx, res, doc, vars)
}

func buildReflectiveOperation(typ gql.OperationType, sel interface{}, vars map[string]interface{}) (res interface{}, doc string) {
	res = &encoding.SelectionUnmarshaler{
		Selection: sel,
	}
	doc = encoding.MustMarshalOperation(&gql.Operation{
		OperationDefinition: gql.OperationDefinition{
			Type:      typ,
			Selection: sel,
		},
		Variables: vars,
	})
	return
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
		"arguments": JSONObject(arguments),
	})
	return m.Job.ID, err
}
