package process

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/core"
	"github.com/deref/exo/jsonutil"
	logcol "github.com/deref/exo/logcol/api"
)

func (provider *Provider) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	// Ensure top-level var directory.
	err := os.Mkdir(provider.VarDir, 0700)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("creating var directory: %w", err)
	}

	// Create var directory for the new process.
	procDir := filepath.Join(provider.VarDir, input.ID)
	if err := os.Mkdir(procDir, 0700); err != nil {
		return nil, fmt.Errorf("creating proc directory: %w", err)
	}

	// Processes are started by default.
	state, err := provider.start(ctx, procDir, input.Spec)
	if err != nil {
		return nil, err
	}

	// Register logs.
	// TODO: Don't do this synchronously here. Use some kind of component hierarchy mechanism.
	collector := log.CurrentLogCollector(ctx)
	for _, role := range []string{"out", "err"} {
		_, err := collector.AddLog(ctx, &logcol.AddLogInput{
			Name:   fmt.Sprintf("%s:%s", input.ID, role),
			Source: filepath.Join(procDir, role),
		})
		if err != nil {
			return nil, fmt.Errorf("adding std%s log: %w", role, err)
		}
	}

	var output core.InitializeOutput
	output.State = jsonutil.MustMarshalString(state)
	return &output, nil
}

func readLine(r io.Reader) (string, error) {
	b := bufio.NewReader(r)
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

func (provider *Provider) Refresh(context.Context, *core.RefreshInput) (*core.RefreshOutput, error) {
	panic("TODO: refresh")
}

func (provider *Provider) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	var state state
	if err := jsonutil.UnmarshalString(input.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	provider.stop(state.Pid)

	// Deregister log streams.
	// TODO: Don't do this synchronously here. Use some kind of component hierarchy mechanism.
	collector := log.CurrentLogCollector(ctx)
	for _, role := range []string{"out", "err"} {
		_, err := collector.RemoveLog(ctx, &logcol.RemoveLogInput{
			Name: fmt.Sprintf("%s:%s", input.ID, role),
		})
		if err != nil {
			return nil, fmt.Errorf("removing std%s log: %w", role, err)
		}
	}

	// Delete var directory.
	procDir := filepath.Join(provider.VarDir, input.ID)
	if err := os.RemoveAll(procDir); err != nil {
		return nil, fmt.Errorf("removing var directory: %w", err)
	}

	return &core.DisposeOutput{State: input.State}, nil
}
