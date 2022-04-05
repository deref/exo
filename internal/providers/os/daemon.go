package os

import (
	"context"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/sdk"
)

type DaemonController struct {
	sdk.PureComponentController[DaemonModel]
}

func NewDaemonController(svc api.Service) sdk.AComponentController {
	return sdk.NewComponentControllerAdapater[DaemonModel](&DaemonController{})
}

type DaemonModel struct {
	DaemonSpec
}

type DaemonSpec struct {
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments,omitempty"`
	Directory   string            `json:"directory,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

func (ctrl *DaemonController) RenderComponent(ctx context.Context, cfg *sdk.ComponentConfig, m *DaemonModel) ([]sdk.RenderedComponent, error) {
	var children []sdk.RenderedComponent
	if cfg.Run {
		// XXX Run supervisor to get logs.
		children = append(children, sdk.RenderedComponent{
			Type: "process",
			Name: "process",
			Spec: ProcessSpec{
				Program:     m.Program,
				Arguments:   m.Arguments,
				Environment: m.Environment,
				Directory:   m.Directory,
			},
		})
	}
	return children, nil
}
