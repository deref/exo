package process

import "context"

type Process interface {
	Start(context.Context, *StartInput) (*StartOutput, error)
	Stop(context.Context, *StopInput) (*StopOutput, error)
}

type StartInput struct {
	State map[string]interface{} `json:"state"`
}

type StartOutput struct {
	State map[string]interface{} `json:"state"`
}

type StopInput struct {
	State map[string]interface{} `json:"state"`
}

type StopOutput struct {
	State map[string]interface{} `json:"state"`
}
