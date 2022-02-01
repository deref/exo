package daemon

import (
	"context"
	"errors"

	"github.com/deref/exo/internal/providers/sdk"
)

type Spec struct {
	Directory       string            `json:"directory,omitempty"`
	Program         string            `json:"program"`
	Arguments       []string          `json:"arguments,omitempty"`
	Environment     map[string]string `json:"environment,omitempty"`
	ShutdownTimeout *int              `json:"shutdownTimeout,omitempty"`
}

type Controller struct{}

func (c *Controller) Render(ctx context.Context, spec *Spec) (*sdk.RenderResult, error) {
	return nil, errors.New("TODO: Render daemon")
}
