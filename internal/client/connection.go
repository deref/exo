package client

import (
	"context"

	"github.com/deref/exo/internal/api"
)

type Connection struct {
	Client
}

func (conn *Connection) Shutdown(ctx context.Context) error {
	// No-op.
	return nil
}

var _ api.Service = (*Connection)(nil)
