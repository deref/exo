package process

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/osutil"
)

var _ core.Lifecycle = (*Process)(nil)

func (p *Process) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	return &core.DependenciesOutput{Components: []string{}}, nil
}

func (p *Process) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	var spec Spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	// Resolve spec into state.
	p.State.Directory = spec.Directory
	p.State.Program = spec.Program
	p.State.Arguments = spec.Arguments
	p.State.Environment = spec.Environment
	p.State.ShutdownGracePeriodSeconds = spec.ShutdownGracePeriodSeconds

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
	// NOTE [PROCESS_STATE_MIGRATION]: This migration copies extra data from spec
	// in to state. After sufficient time has passed from October 2021, this and
	// the notes referencing it can be removed.
	var spec Spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}
	p.State.Directory = spec.Directory
	p.State.Program = spec.Program
	p.State.Arguments = spec.Arguments
	p.State.Environment = spec.Environment
	p.State.ShutdownGracePeriodSeconds = spec.ShutdownGracePeriodSeconds

	p.refresh()
	return &core.RefreshOutput{}, nil
}

func (p *Process) refresh() {
	if osutil.IsValidPid(p.SupervisorPid) && osutil.IsValidPid(p.Pid) {
		return
	}
	p.State.reset()
}

func (p *Process) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	if err := p.stop(nil); err != nil {
		return nil, err
	}
	return &core.DisposeOutput{}, nil
}
