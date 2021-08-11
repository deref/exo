package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"syscall"

	"github.com/deref/exo"
	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/task"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/upgrade"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
)

type Kernel struct {
	VarDir      string
	Store       state.Store
	TaskTracker *task.TaskTracker
}

func (kern *Kernel) CreateWorkspace(ctx context.Context, input *api.CreateWorkspaceInput) (*api.CreateWorkspaceOutput, error) {
	id := gensym.RandomBase32()
	_, err := kern.Store.AddWorkspace(ctx, &state.AddWorkspaceInput{
		ID:   id,
		Root: input.Root,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateWorkspaceOutput{
		ID: id,
	}, nil
}

func (kern *Kernel) DescribeWorkspaces(ctx context.Context, input *api.DescribeWorkspacesInput) (*api.DescribeWorkspacesOutput, error) {
	output, err := kern.Store.DescribeWorkspaces(ctx, &state.DescribeWorkspacesInput{})
	if err != nil {
		return nil, err
	}
	workspaces := make([]api.WorkspaceDescription, len(output.Workspaces))
	for i, workspace := range output.Workspaces {
		workspaces[i] = api.WorkspaceDescription{
			ID:   workspace.ID,
			Root: workspace.Root,
		}
	}
	return &api.DescribeWorkspacesOutput{
		Workspaces: workspaces,
	}, nil
}

func (kern *Kernel) FindWorkspace(ctx context.Context, input *api.FindWorkspaceInput) (*api.FindWorkspaceOutput, error) {
	output, err := kern.Store.FindWorkspace(ctx, &state.FindWorkspaceInput{
		Path: input.Path,
	})
	if err != nil {
		return nil, err
	}
	return &api.FindWorkspaceOutput{
		ID: output.ID,
	}, nil
}

func (kern *Kernel) GetVersion(ctx context.Context, input *api.GetVersionInput) (*api.GetVersionOutput, error) {
	tel := telemetry.FromContext(ctx)
	installed := exo.Version
	current := true
	var latest *string
	if tel.IsEnabled() {
		latestVersion, err := tel.LatestVersion(ctx)
		if err != nil {
			return nil, err
		}
		latest = &latestVersion
		current = installed >= latestVersion
	}

	return &api.GetVersionOutput{
		Installed: installed,
		Latest:    latest,
		Current:   current,
	}, nil
}

func (kern *Kernel) Upgrade(ctx context.Context, input *api.UpgradeInput) (*api.UpgradeOutput, error) {
	if upgrade.IsManaged {
		return nil, errutil.WithHTTPStatus(http.StatusBadRequest, errors.New("exo installed with system package manager"))
	}
	err := upgrade.UpgradeSelf()
	if err != nil {
		return nil, err
	}
	defer restart(ctx)

	return &api.UpgradeOutput{}, nil
}

func (kern *Kernel) Panic(ctx context.Context, input *api.PanicInput) (*api.PanicOutput, error) {
	message := input.Message
	if input.Message == "" {
		message = "test error"
	}
	panic(message)
}

func (kern *Kernel) Ping(context.Context, *api.PingInput) (*api.PingOutput, error) {
	return &api.PingOutput{}, nil
}

// restart replaces the current process with a new copy of itself. This is useful
// after downloading a new binary.
func restart(ctx context.Context) {
	// TODO: Shutdown gracefully.

	// Replace the current process with the newer version of itself.
	exitWithError := func(err error) {
		fmt.Printf("Fatal error restarting exo: %v\n", err)
		os.Exit(1)
	}

	cmd, err := os.Executable()
	if err != nil {
		exitWithError(fmt.Errorf("getting current executable: %w", err))
	}

	// Since the exo process is likely a specific version that `exo` is linked to,
	// we check to see if there is an `exo` symlink in the same directory as the
	// current executable, and if so, we run that instead.
	dir := path.Dir(cmd)
	symlinkPath := path.Join(dir, "exo")
	if isSymlink, _ := osutil.IsSymlink(symlinkPath); isSymlink {
		dest, err := os.Readlink(symlinkPath)
		if err != nil {
			exitWithError(fmt.Errorf("following exo symlink: %w", err))
		}
		if !path.IsAbs(dest) {
			dest = path.Join(dir, dest)
		}

		cmd = path.Clean(dest)
	}

	if err := syscall.Exec(cmd, append([]string{cmd}, os.Args[1:]...), os.Environ()); err != nil {
		exitWithError(fmt.Errorf("forking new exo process at %s: %w", cmd, err))
	}
}

func (kern *Kernel) DescribeTasks(ctx context.Context, input *api.DescribeTasksInput) (*api.DescribeTasksOutput, error) {
	underlying, err := kern.TaskTracker.Store.DescribeTasks(ctx, &taskapi.DescribeTasksInput{
		JobIDs: input.JobIDs,
	})
	if err != nil {
		return nil, err
	}
	var output api.DescribeTasksOutput
	output.Tasks = make([]api.TaskDescription, len(underlying.Tasks))
	for i, t := range underlying.Tasks {
		output.Tasks[i] = api.TaskDescription{
			ID:       t.ID,
			JobID:    t.JobID,
			ParentID: t.ParentID,
			Name:     t.Name,
			Status:   t.Status,
			Message:  t.Message,
			Created:  t.Created,
			Updated:  t.Updated,
			Started:  t.Started,
			Finished: t.Finished,
		}
	}
	return &output, nil
}
