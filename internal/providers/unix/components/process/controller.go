// TODO: Generate these.

package process

import (
	"fmt"

	"github.com/deref/util-go/jsonutil"
)

func (p *Process) InitResource() error {
	if err := jsonutil.UnmarshalString(p.ComponentSpec, &p.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(p.ComponentState, &p.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (p *Process) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(p.State)
}
