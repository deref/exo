package volume

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"gopkg.in/yaml.v3"
)

func (v *Volume) InitResource() error {
	if err := yaml.Unmarshal([]byte(v.ComponentSpec), &v.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(v.ComponentState, &v.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (v *Volume) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(v.State)
}
