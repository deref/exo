// TODO: Generate these.

package process

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (p *Process) InitResource(componentID, spec, state string) error {
	p.ComponentID = componentID
	if err := jsonutil.UnmarshalString(spec, &p.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &p.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (p *Process) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(p.State)
}
