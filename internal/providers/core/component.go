package core

import "github.com/deref/exo/internal/util/logging"

type Component struct {
	ComponentID   string
	WorkspaceRoot string
	Logger        logging.Logger
}
