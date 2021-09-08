// TODO: Generate these.

package process

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/goccy/go-yaml"
)

func (p *Process) InitResource() error {
	if err := yaml.Unmarshal([]byte(p.ComponentSpec), &p.Spec); err != nil {
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
