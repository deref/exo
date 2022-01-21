package resolvers

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/providers/os/resources/file"
	"github.com/deref/exo/internal/providers/os/resources/process"
	"github.com/deref/exo/internal/providers/sdk"
)

func getResourceController(ctx context.Context, typ string) (*sdk.Controller, error) {
	impl := getResourceControllerImpl(ctx, typ)
	if impl == nil {
		return nil, fmt.Errorf("no resource controller for type: %s", typ)
	}
	return sdk.NewController(impl), nil
}

// TODO: Dynamic registry with qualified type identifiers.
func getResourceControllerImpl(ctx context.Context, typ string) interface{} {
	switch typ {
	case "file":
		return &file.Controller{}
	case "process":
		return &process.Controller{}
	default:
		return nil
	}
}
