package resolvers

import (
	"context"

	"github.com/deref/exo/internal/providers/os/daemon"
	"github.com/deref/exo/internal/providers/os/file"
	"github.com/deref/exo/internal/providers/os/process"
	"github.com/deref/exo/internal/providers/sdk"
)

// TODO: Dynamic registry with qualified type identifiers.

func isProcessController(typ string) bool {
	switch typ {
	case "daemon", "process", "container":
		return true
	default:
		return false
	}
}

func getController(ctx context.Context, typ string) *sdk.Controller {
	impl := getControllerImpl(ctx, typ)
	if impl == nil {
		return nil
	}
	return sdk.NewController(impl)
}

func getControllerImpl(ctx context.Context, typ string) any {
	switch typ {
	case "daemon":
		return &daemon.Controller{}
	case "file":
		return &file.Controller{}
	case "process":
		return &process.Controller{}
	default:
		return nil
	}
}
