package volume

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/goccy/go-yaml"
)

func (v *Volume) InitResource(spec, state string) error {
	if err := yaml.Unmarshal([]byte(spec), &v.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &v.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (v *Volume) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(v.State)
}
