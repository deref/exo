package process

import (
	"context"

	"github.com/deref/exo/internal/util/jsonutil"
)

type Model struct {
	HostID string `json:"hostId"`
	Pid    *int   `json:"pid"`
	// Absolute path to program to run on host.
	Program     string            `json:"program"`
	Arguments   []string          `json:"arguments"`
	Directory   string            `json:"directory"`
	Environment map[string]string `json:"environment"`
}

func (m *Model) UnmarshalModel(ctx context.Context, s string) error {
	return jsonutil.UnmarshalString(s, m)
}

func (m *Model) MarshalModel(ctx context.Context) (string, error) {
	return jsonutil.MarshalString(m)
}
