package process

import (
	"bufio"
	"context"
	"errors"
	"io"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/osutil"
)

func (p *Process) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	// Processes are started by default.
	if err := p.start(ctx); err != nil {
		return nil, err
	}
	return &core.InitializeOutput{}, nil
}

func readLine(r io.Reader) (string, error) {
	b := bufio.NewReaderSize(r, 4096)
	line, isPrefix, err := b.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("line too long")
	}
	return string(line), nil
}

func (p *Process) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	p.refresh()
	return &core.RefreshOutput{}, nil
}

func (p *Process) refresh() {
	if !osutil.IsValidPid(p.SupervisorPid) {
		p.State.clear()
	}
}

func (p *Process) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	p.stop(false)
	return &core.DisposeOutput{}, nil
}
