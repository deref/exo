package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/deref/graphql-go/encoding"
	"github.com/deref/graphql-go/gql"
)

type Service interface {
	Shutdown(context.Context) error
	// Execute a GraphQL operation synchronously, decoding the response data in to res.
	// Implementations should also respect CurrentContextVariables.
	Do(ctx context.Context, res interface{}, doc string, vars map[string]interface{}) error
	// Begin a GraphQL subscription operation. Decodes responses in to the result
	// of invoking newRes for each event response. Decoded responses are
	// delivered to to consumers via Subscription.Events().
	Subscribe(ctx context.Context, newRes func() interface{}, doc string, vars map[string]interface{}) Subscription
}

type Subscription interface {
	// Yields events until the subscription is stopped.  The element type is
	// specified via the original Subscribe() method call.
	Events() <-chan interface{}
	Err() error
	Stop()
}

func Query(ctx context.Context, svc Service, q interface{}, vars map[string]interface{}) error {
	res := newSelectionUnmarshaler(q)
	doc := marshalReflectiveOperation(gql.Query, q, vars)
	return svc.Do(ctx, res, doc, vars)
}

func Mutate(ctx context.Context, svc Service, m interface{}, vars map[string]interface{}) error {
	res := newSelectionUnmarshaler(m)
	doc := marshalReflectiveOperation(gql.Mutation, m, vars)
	return svc.Do(ctx, res, doc, vars)
}

// Like Query and Mutate, but the reflective structure will not be modified directly.
// Instead, a new instance will be allocated for each event.
func Subscribe(ctx context.Context, svc Service, s interface{}, vars map[string]interface{}) Subscription {
	resType := reflect.TypeOf(s).Elem()
	newRes := func() interface{} {
		return newSelectionUnmarshaler(reflect.New(resType).Interface())
	}
	doc := marshalReflectiveOperation(gql.Subscription, s, vars)
	return svc.Subscribe(ctx, newRes, doc, vars)
}

func newSelectionUnmarshaler(sel interface{}) *encoding.SelectionUnmarshaler {
	return &encoding.SelectionUnmarshaler{
		Selection: sel,
	}
}

func marshalReflectiveOperation(typ gql.OperationType, sel interface{}, vars map[string]interface{}) string {
	return encoding.MustMarshalOperation(&gql.Operation{
		OperationDefinition: gql.OperationDefinition{
			Type:      typ,
			Selection: sel,
		},
		Variables: vars,
	})
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

// Extracts the operation data from a response structures. The response data
// must contain exactly one top-level field.
func OperationData(resp interface{}) interface{} {
	switch resp := resp.(type) {
	case *encoding.SelectionUnmarshaler:
		return OperationData(resp.Selection)

	case map[string]interface{}:
		if len(resp) != 1 {
			panic(fmt.Errorf("expected map have exactly one entry, found %d", len(resp)))
		}
		for _, v := range resp {
			return v
		}
		panic("unreachable")

	default:
		v := reflect.ValueOf(resp)
		if v.Kind() == reflect.Ptr {
			return OperationData(v.Elem().Interface())
		}
		if v.Kind() != reflect.Struct {
			panic(fmt.Errorf("cannot extract operation data from %T", resp))
		}
		if v.NumField() != 1 {
			panic(fmt.Errorf("expected %T to have exactly one field, found %d", resp, v.NumField()))
		}
		return v.Field(0).Interface()
	}
}
