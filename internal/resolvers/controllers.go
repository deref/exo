package resolvers

import (
	"context"

	"github.com/deref/exo/internal/providers/os/components/daemon"
	"github.com/deref/exo/internal/providers/os/resources/file"
	"github.com/deref/exo/internal/providers/os/resources/process"
	"github.com/deref/exo/internal/providers/sdk"
)

// TODO: Dynamic registry with qualified type identifiers.

func getResourceController(ctx context.Context, typ string) *sdk.ResourceController {
	impl := getResourceControllerImpl(ctx, typ)
	if impl == nil {
		return nil
	}
	return sdk.NewResourceController(impl)
}

func getComponentController(ctx context.Context, typ string) *sdk.ComponentController {
	impl := getComponentControllerImpl(ctx, typ)
	if impl == nil {
		return nil
	}
	return sdk.NewComponentController(impl)
}

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

func getComponentControllerImpl(ctx context.Context, typ string) interface{} {
	switch typ {
	case "daemon":
		return &daemon.Controller{}
	default:
		return nil
	}
}
