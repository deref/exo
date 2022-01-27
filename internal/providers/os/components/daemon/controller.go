package daemon

import (
	"context"
	"errors"

	"github.com/deref/exo/internal/providers/sdk"
)

type Specification struct{}

type Controller struct{}

func (c *Controller) Render(ctx context.Context, spec Specification) (*sdk.RenderResult, error) {
	return nil, errors.New("TODO: Render daemon")
}
