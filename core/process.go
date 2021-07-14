package core

import "context"

type Process interface {
	Start(context.Context, *StartInput) (*StartOutput, error)
	Stop(context.Context, *StopInput) (*StopOutput, error)
}

type StartInput struct {
	ID    string `json:"id"`
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type StartOutput struct {
	State string `json:"state"`
}

type StopInput struct {
	ID    string `json:"id"`
	Spec  string `json:"spec"`
	State string `json:"state"`
}

type StopOutput struct {
	State string `json:"state"`
}
