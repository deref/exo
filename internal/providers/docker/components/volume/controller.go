package volume

import (
	"fmt"

	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
)

func (v *Volume) InitResource(componentID, spec, state string) error {
	v.ComponentID = componentID
	if err := docker.LoadSpec(spec, &v.Spec, v.ComponentBase.WorkspaceEnvironment); err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &v.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (v *Volume) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(v.State)
}
