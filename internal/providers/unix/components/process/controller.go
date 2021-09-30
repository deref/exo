// TODO: Generate these.

package process

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (p *Process) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(p.ComponentState, &p.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (p *Process) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(p.State)
}
