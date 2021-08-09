package core

import "github.com/deref/exo/internal/util/logging"

type Component interface {
	GetComponentID() string
}

type ComponentBase struct {
	ComponentID   string
	WorkspaceRoot string
	Logger        logging.Logger
}

func (base *ComponentBase) GetComponentID() string {
	return base.ComponentID
}
