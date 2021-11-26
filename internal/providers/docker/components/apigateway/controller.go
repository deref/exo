package apigateway

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (ag *APIGateway) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(ag.ComponentState, &ag.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (ag *APIGateway) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(ag.State)
}
