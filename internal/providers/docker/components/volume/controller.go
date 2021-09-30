package volume

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (v *Volume) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(v.ComponentState, &v.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (v *Volume) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(v.State)
}
