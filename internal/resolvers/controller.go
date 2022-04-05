package resolvers

import (
	"context"

	"github.com/deref/exo/internal/providers/os"
	"github.com/deref/exo/sdk"
)

// TODO: Dynamic registry with qualified type identifiers.

func (r *QueryResolver) componentControllerByType(ctx context.Context, typ string) sdk.AComponentController {
	switch typ {
	case "daemon":
		return os.NewDaemonController(r.Service)
	default:
		return r.resourceControllerByType(ctx, typ)
	}
}

func (r *QueryResolver) resourceControllerByType(ctx context.Context, typ string) *sdk.ResourceComponentController {
	switch typ {
	case "file":
		return os.NewFileController(r.Service)
	case "process":
		return os.NewProcessController(r.Service)
	default:
		return nil
	}
}
