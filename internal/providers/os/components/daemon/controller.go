package daemon

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/providers/sdk"
)

type Component struct {
	sdk.ComponentConfig
	Spec Spec `json:"spec"`
}

type Spec struct {
	Directory       string            `json:"directory,omitempty"`
	Program         string            `json:"program"`
	Arguments       []string          `json:"arguments,omitempty"`
	Environment     map[string]string `json:"environment,omitempty"`
	ShutdownTimeout *int              `json:"shutdownTimeout,omitempty"`
}

type Controller struct{}

func (ctrl *Controller) Render(ctx context.Context, component *Component) (*sdk.RenderResult, error) {
	return nil, fmt.Errorf("TODO: Render daemon: %#v", component)
}
