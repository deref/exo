package core

import (
	"github.com/deref/exo/internal/util/logging"
)

type Component interface {
	GetComponentID() string

	// TODO: Rethink these.
	IsDeleted() bool
	MarkDeleted()
}

type ComponentBase struct {
	ComponentID   string
	WorkspaceRoot string
	Logger        logging.Logger
	isDeleted     bool
}

func (c *ComponentBase) GetComponentID() string {
	return c.ComponentID
}

func (c *ComponentBase) IsDeleted() bool {
	return c.isDeleted
}

func (c *ComponentBase) MarkDeleted() {
	c.isDeleted = true
}
