package daemon

import (
	"context"

	"github.com/deref/exo/internal/providers/os/process"
	"github.com/deref/exo/internal/providers/sdk"
	"github.com/deref/exo/internal/util/logging"
)

func (ctrl *Controller) Render(ctx context.Context, component *Component) ([]sdk.RenderedComponent, error) {
	logging.Infof(ctx, "XXX Render daemon: %#v", component)
	spec := component.Spec
	var children []sdk.RenderedComponent
	if component.Run {
		// XXX Run supervisor to get logs.
		children = append(children, sdk.RenderedComponent{
			Type: "process",
			Name: "process",
			Spec: process.Spec{
				Program:     spec.Program,
				Arguments:   spec.Arguments,
				Environment: spec.Environment,
				Directory:   spec.Directory,
			},
		})
	}
	return children, nil
}
