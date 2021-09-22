package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/deref/exo"
	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/task"
	taskapi "github.com/deref/exo/internal/task/api"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/template"
	"github.com/deref/exo/internal/upgrade"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/otiai10/copy"
)

type Kernel struct {
	VarDir      string
	Store       state.Store
	TaskTracker *task.TaskTracker
}

func (kern *Kernel) ListTemplates(ctx context.Context, input *api.ListTemplatesInput) (*api.ListTemplatesOutput, error) {
	// This list is the set of directories here: https://github.com/railwayapp/starters/tree/master/examples
	// TODO: make this dynamic
	return &api.ListTemplatesOutput{
		TemplateNames: []string{
			"ackee",
			"blitzjs",
			"calendso",
			"code-server",
			"cusdis",
			"deno",
			"discordjs-typescript",
			"discordjs",
			"discordpy",
			"djangopy",
			"elixir-phoenix",
			"eris",
			"expressjs-mongoose",
			"expressjs-postgres",
			"expressjs-prisma",
			"expressjs",
			"fastapi",
			"fastify",
			"flask",
			"ghost",
			"gin",
			"hapi",
			"hasura",
			"http-nodejs",
			"laravel",
			"n8n",
			"next-notion-blog",
			"nextjs-auth-mongo",
			"nextjs-prisma",
			"nocodb",
			"node_modules",
			"nuxtjs",
			"rails-blog",
			"rails-starter",
			"rocket",
			"rust-wasm",
			"shiori",
			"sinatra",
			"slack-webhook",
			"starlette",
			"strapi",
			"svelte-kit",
			"telegram-bot",
			"umami",
		},
	}, nil
}

func (kern *Kernel) CreateProject(ctx context.Context, input *api.CreateProjectInput) (*api.CreateProjectOutput, error) {
	projectDir := input.Root
	if !path.IsAbs(projectDir) {
		return &api.CreateProjectOutput{}, errors.New("path must be absolute")
	}

	var templateDir string
	if input.TemplateName != nil {
		var err error
		templateDir, err = template.GetTemplateFiles(ctx, *input.TemplateName)
		if err != nil {
			return &api.CreateProjectOutput{}, fmt.Errorf("getting template files: %w", err)
		}
	}

	err := os.Mkdir(projectDir, 0750)
	if err != nil {
		return &api.CreateProjectOutput{}, fmt.Errorf("making project dir: %w", err)
	}

	if templateDir != "" {
		copy.Copy(templateDir, projectDir, copy.Options{
			OnSymlink: func(string) copy.SymlinkAction {
				return copy.Deep
			},
		})
	}

	result, err := kern.CreateWorkspace(ctx, &api.CreateWorkspaceInput{Root: projectDir})
	if err != nil {
		return &api.CreateProjectOutput{}, fmt.Errorf("creating workspace: %w", err)
	}
	return &api.CreateProjectOutput{WorkspaceID: result.ID}, err
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
			ID:          workspace.ID,
			Root:        workspace.Root,
			DisplayName: workspace.DisplayName,
		}
	}
	return &api.DescribeWorkspacesOutput{
		Workspaces: workspaces,
	}, nil
}

func (kern *Kernel) ResolveWorkspace(ctx context.Context, input *api.ResolveWorkspaceInput) (*api.ResolveWorkspaceOutput, error) {
	output, err := kern.Store.ResolveWorkspace(ctx, &state.ResolveWorkspaceInput{
		Ref: input.Ref,
	})
	if err != nil {
		return nil, err
	}
	return &api.ResolveWorkspaceOutput{
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

func (kern *Kernel) Exit(context.Context, *api.ExitInput) (*api.ExitOutput, error) {
	// Return immediately, shutdown asynchronously.
	go func() {
		ownPid := os.Getpid()
		_ = osutil.TerminateProcessWithTimeout(ownPid, 5*time.Second)
	}()

	return &api.ExitOutput{}, nil
}

// restart replaces the current process with a new copy of itself. This is useful
// after downloading a new binary.
func restart(ctx context.Context) {
	// TODO: Shutdown gracefully. Since shutdown code is tightly tied to process termination,
	// a refactor is probably needed so that we have a single place to hook into for shutdown
	// code.

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
	for i, t1 := range underlying.Tasks {
		t2 := api.TaskDescription{
			ID:       t1.ID,
			JobID:    t1.JobID,
			ParentID: t1.ParentID,
			Name:     t1.Name,
			Status:   t1.Status,
			Message:  t1.Message,
			Created:  t1.Created,
			Updated:  t1.Updated,
			Started:  t1.Started,
			Finished: t1.Finished,
		}
		if t1.Progress != nil {
			t2.Progress = &api.TaskProgress{
				Current: t1.Progress.Current,
				Total:   t1.Progress.Total,
			}
		}
		output.Tasks[i] = t2
	}
	return &output, nil
}
