package apigateway

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (ag *APIGateway) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(ag.ComponentState, &ag.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	ag.Container.State = ag.State.State
	return nil
}

func (ag *APIGateway) MarshalState() (state string, err error) {
	ag.State.State = ag.Container.State
	return jsonutil.MarshalString(ag.State)
}
