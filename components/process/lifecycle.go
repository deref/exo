package process

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	core "github.com/deref/exo/core/api"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/osutil"
)

func (provider *Provider) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	// Processes are started by default.
	state, err := provider.start(ctx, input.ID, input.Spec)
	if err != nil {
		return nil, err
	}

	var output core.InitializeOutput
	output.State = jsonutil.MustMarshalString(state)
	return &output, nil
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

func (provider *Provider) Update(context.Context, *core.UpdateInput) (*core.UpdateOutput, error) {
	panic("TODO: update")
}

func (provider *Provider) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	var state State
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	provider.refresh(&state)

	var output core.RefreshOutput
	output.State = jsonutil.MustMarshalString(state)
	return &output, nil
}

func (provider *Provider) refresh(state *State) {
	if !osutil.IsValidPid(state.Pid) {
		state.Pid = 0
	}
}

func (provider *Provider) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	var state State
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	provider.stop(state.Pid)

	return &core.DisposeOutput{State: input.State}, nil
}
