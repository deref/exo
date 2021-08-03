package core

import "github.com/deref/exo/util/logging"

type Component struct {
	ComponentID   string
	WorkspaceRoot string
	Logger        logging.Logger
}
