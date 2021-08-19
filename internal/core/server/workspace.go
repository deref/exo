package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/gensym"
	logd "github.com/deref/exo/internal/logd/api"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/deref/exo/internal/providers/core"
	"github.com/deref/exo/internal/providers/core/components/invalid"
	"github.com/deref/exo/internal/providers/core/components/log"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/components/container"
	"github.com/deref/exo/internal/providers/docker/components/network"
	"github.com/deref/exo/internal/providers/docker/components/volume"
	"github.com/deref/exo/internal/providers/unix/components/process"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/exo/internal/util/pathutil"
	dockerclient "github.com/docker/docker/client"
	psprocess "github.com/shirou/gopsutil/v3/process"
)

type Workspace struct {
	ID          string
	VarDir      string
	Store       state.Store
	SyslogPort  uint
	Logger      logging.Logger // TODO: Embed in context, so it can be annotated with request info.
	Docker      *dockerclient.Client
	TaskTracker *task.TaskTracker
}

func (ws *Workspace) Describe(ctx context.Context, input *api.DescribeInput) (*api.DescribeOutput, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return nil, err
	}
	return &api.DescribeOutput{
		Description: *description,
	}, nil
}

func (ws *Workspace) describe(ctx context.Context) (*api.WorkspaceDescription, error) {
	output, err := ws.Store.DescribeWorkspaces(ctx, &state.DescribeWorkspacesInput{
		IDs: []string{ws.ID},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Workspaces) != 1 {
		return nil, fmt.Errorf("invalid workspace: %q", ws.ID)
	}
	return &api.WorkspaceDescription{
		ID:   ws.ID,
		Root: output.Workspaces[0].Root,
	}, nil
}

func (ws *Workspace) Destroy(ctx context.Context, input *api.DestroyInput) (*api.DestroyOutput, error) {
	job := ws.TaskTracker.StartTask(ctx, "destroying")
	go func() {
		defer job.Finish()
		filter := componentFilter{}
		ws.goControlComponents(job, filter, func(ctx context.Context, lifecycle api.Lifecycle) error {
			return ws.deleteComponent(ctx, lifecycle)
		})
		if err := job.Wait(); err != nil {
			return
		}
		if _, err := ws.Store.RemoveWorkspace(ctx, &state.RemoveWorkspaceInput{
			ID: ws.ID,
		}); err != nil {
			job.Fail(fmt.Errorf("removing workspace from store: %w", err))
			return
		}
	}()
	return &api.DestroyOutput{
		JobID: job.JobID(),
	}, nil
}

func (ws *Workspace) Apply(ctx context.Context, input *api.ApplyInput) (*api.ApplyOutput, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return nil, fmt.Errorf("describing workspace: %w", err)
	}
	res := ws.loadManifest(description.Root, input)
	if res.Err != nil {
		return nil, res.Err
	}
	m := res.Manifest

	describeOutput, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	// Index old components by name.
	oldComponents := make(map[string]api.ComponentDescription, len(describeOutput.Components))
	for _, component := range describeOutput.Components {
		oldComponents[component.Name] = component
	}

	// TODO: Handle partial failures.

	job := ws.TaskTracker.StartTask(ctx, "applying")
	go func() {
		defer job.Finish()

		// Apply component upserts.
		newComponents := make(map[string]manifest.Component, len(m.Components))
		for _, newComponent := range m.Components {
			newComponent := newComponent
			name := newComponent.Name
			newComponents[name] = newComponent
			if oldComponent, exists := oldComponents[name]; exists {
				// Update existing component.
				job.Go("updating "+name, func(t *task.Task) error {
					return ws.updateComponent(t, oldComponent, newComponent)
				})
			} else {
				// Create new component.
				job.Go("adding "+name, func(t *task.Task) error {
					_, err := ws.createComponent(t, newComponent)
					return err
				})
			}
		}

		// Apply component deletions.
		for name, oldComponent := range oldComponents {
			name := name
			oldComponent := oldComponent
			if _, keep := newComponents[name]; keep {
				continue
			}
			job.Go("deleting "+name, func(*task.Task) error {
				return ws.control(job.Context, oldComponent, func(ctx context.Context, lifecycle api.Lifecycle) error {
					return ws.deleteComponent(job.Context, lifecycle)
				})
			})
		}

	}()

	return &api.ApplyOutput{
		Warnings: res.Warnings,
		JobID:    job.ID(),
	}, nil
}

func (ws *Workspace) Resolve(ctx context.Context, input *api.ResolveInput) (*api.ResolveOutput, error) {
	storeOutput, err := ws.Store.Resolve(ctx, &state.ResolveInput{
		WorkspaceID: ws.ID,
		Refs:        input.Refs,
	})
	if err != nil {
		return nil, err
	}
	var output api.ResolveOutput
	output.IDs = make([]*string, len(storeOutput.IDs))
	for i, id := range storeOutput.IDs {
		output.IDs[i] = id
	}
	return &output, err
}

func (ws *Workspace) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (*api.DescribeComponentsOutput, error) {
	stateOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		IDs:         input.IDs,
		Types:       input.Types,
	})
	if err != nil {
		return nil, err
	}
	output := &api.DescribeComponentsOutput{
		Components: []api.ComponentDescription{},
	}
	for _, component := range stateOutput.Components {
		output.Components = append(output.Components, api.ComponentDescription{
			ID:          component.ID,
			Name:        component.Name,
			Type:        component.Type,
			Spec:        component.Spec,
			State:       component.State,
			Created:     component.Created,
			Initialized: component.Initialized,
			Disposed:    component.Disposed,
		})
	}
	return output, nil
}

func (ws *Workspace) newController(ctx context.Context, desc api.ComponentDescription) Controller {
	description, err := ws.describe(ctx)
	if err != nil {
		return &invalid.Invalid{
			Err: fmt.Errorf("workspace error: %w", err),
		}
	}
	base := core.ComponentBase{
		ComponentID:          desc.ID,
		ComponentName:        desc.Name,
		ComponentSpec:        desc.Spec,
		ComponentState:       desc.State,
		WorkspaceRoot:        description.Root,
		WorkspaceEnvironment: ws.getEnvironment(),
		Logger:               ws.Logger,
	}
	switch desc.Type {
	case "process":
		return &process.Process{
			ComponentBase: base,
			SyslogPort:    ws.SyslogPort,
		}

	case "container":
		return &container.Container{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
			SyslogPort: ws.SyslogPort,
		}

	case "network":
		return &network.Network{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
		}

	case "volume":
		return &volume.Volume{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
		}

	default:
		return &invalid.Invalid{
			Err: fmt.Errorf("unsupported component type: %q", desc.Type),
		}
	}
}

// TODO: Use workspace-defined environments, rather than ambient unix environment.
func (ws *Workspace) getEnvironment() map[string]string {
	env := make(map[string]string)
	for _, assign := range os.Environ() {
		parts := strings.SplitN(assign, "=", 2)
		key := parts[0]
		val := parts[1]
		env[key] = val
	}
	return env
}

func (ws *Workspace) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (*api.CreateComponentOutput, error) {
	id, err := ws.createComponent(ctx, manifest.Component{
		Name: input.Name,
		Type: input.Type,
		Spec: input.Spec,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateComponentOutput{
		ID: id,
	}, nil
}

func (ws *Workspace) createComponent(ctx context.Context, component manifest.Component) (id string, err error) {
	if err := manifest.ValidateName(component.Name); err != nil {
		return "", errutil.HTTPErrorf(http.StatusBadRequest, "component name %q invalid: %w", component.Name, err)
	}

	id = gensym.RandomBase32()

	if _, err := ws.Store.AddComponent(ctx, &state.AddComponentInput{
		WorkspaceID: ws.ID,
		ID:          id,
		Name:        component.Name,
		Type:        component.Type,
		Spec:        component.Spec,
		Created:     chrono.NowString(ctx),
	}); err != nil {
		return "", fmt.Errorf("adding component: %w", err)
	}

	if err := ws.control(ctx, api.ComponentDescription{
		// Construct a synthetic component description to avoid re-reading after
		// the add. Only the fields needed by control are included.
		// TODO: Store.AddComponent could return a compponent description?
		ID:   id,
		Type: component.Type,
		Spec: component.Spec,
	}, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Initialize(ctx, &api.InitializeInput{})
		return err
	}); err != nil {
		return "", err
	}

	// XXX this now double-patches the component to set Initialized timestamp. Optimize?
	if _, err := ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:          id,
		Initialized: chrono.NowString(ctx),
	}); err != nil {
		return "", fmt.Errorf("modifying component after initialization: %w", err) // XXX this message seems incorrect.
	}

	return id, nil
}

func (ws *Workspace) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (*api.UpdateComponentOutput, error) {
	panic("TODO: UpdateComponent") // XXX can implement this now.
}

func (ws *Workspace) updateComponent(ctx context.Context, oldComponent api.ComponentDescription, newComponent manifest.Component) error {
	// TODO: Most updates should be accomplished without a full replacement; especially when there are no spec changes!
	if err := ws.control(ctx, oldComponent, func(ctx context.Context, lifecycle api.Lifecycle) error {
		return ws.deleteComponent(ctx, lifecycle)
	}); err != nil {
		return fmt.Errorf("delete %q for replacement: %w", oldComponent.Name, err)
	}
	if _, err := ws.createComponent(ctx, newComponent); err != nil {
		return fmt.Errorf("adding replacement %q: %w", newComponent.Name, err)
	}
	return nil
}

func (ws *Workspace) RefreshComponents(ctx context.Context, input *api.RefreshComponentsInput) (*api.RefreshComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "refreshing", filter, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Refresh(ctx, &api.RefreshInput{})
		return err
	})
	return &api.RefreshComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) DisposeComponents(ctx context.Context, input *api.DisposeComponentsInput) (*api.DisposeComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "disposing", filter, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Dispose(ctx, &api.DisposeInput{})
		return err
	})
	return &api.DisposeComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) resolveRef(ctx context.Context, ref string) (string, error) {
	resolved, err := ws.resolveRefs(ctx, []string{ref})
	if err != nil {
		return "", err
	}
	return resolved[0], nil
}

func (ws *Workspace) resolveRefs(ctx context.Context, refs []string) ([]string, error) {
	resolveOutput, err := ws.Resolve(ctx, &api.ResolveInput{Refs: refs})
	if err != nil {
		return nil, err
	}
	results := make([]string, len(refs))
	for i, id := range resolveOutput.IDs {
		if id == nil {
			return nil, errutil.HTTPErrorf(http.StatusBadRequest, "unresolvable: %q", refs[i])
		}
		results[i] = *id
	}
	return results, nil
}

func (ws *Workspace) DeleteComponents(ctx context.Context, input *api.DeleteComponentsInput) (*api.DeleteComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "deleting", filter, func(ctx context.Context, lifecycle api.Lifecycle) error {
		return ws.deleteComponent(ctx, lifecycle)
	})
	return &api.DeleteComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) deleteComponent(ctx context.Context, lifecycle api.Lifecycle) error {
	_, err := lifecycle.Dispose(ctx, &api.DisposeInput{})
	if err != nil {
		return err
	}
	lifecycle.(core.Component).MarkDeleted()
	return nil
}

func (ws *Workspace) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	components, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{
		Types: processTypes,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	// Find all logs in component hierarchy.
	// TODO: More general handling of log groups, subcomponents, etc.
	var logGroups []string
	var logStreams []string
	streamToGroup := make(map[string]int)
	for _, component := range components.Components {
		// XXX Janky provider inference. See note: [LOG_COMPONENTS].
		var provider string
		switch component.Type {
		case "process":
			provider = "unix"
		case "container":
			provider = "docker"
		}
		for _, streamName := range log.ComponentLogNames(provider, component.ID) {
			streamToGroup[streamName] = len(logGroups)
			logStreams = append(logStreams, streamName)
		}
		logGroups = append(logGroups, component.ID)
	}

	// Initialize output and index by log group name.
	logs := make([]api.LogDescription, len(logGroups))
	for i, logGroup := range logGroups {
		logs[i] = api.LogDescription{
			Name: logGroup,
		}
	}

	// Decorate output with information from the log collector.
	collector := log.CurrentLogCollector(ctx)
	collectorLogs, err := collector.DescribeLogs(ctx, &logd.DescribeLogsInput{
		Names: logStreams,
	})
	if err != nil {
		return nil, err
	}
	for _, collectorLog := range collectorLogs.Logs {
		groupIndex, ok := streamToGroup[collectorLog.Name]
		if !ok {
			continue
		}
		group := &logs[groupIndex]
		group.LastEventAt = combineLastEventAt(group.LastEventAt, collectorLog.LastEventAt)
	}
	return &api.DescribeLogsOutput{Logs: logs}, nil
}

func combineLastEventAt(a, b *string) *string {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if strings.Compare(*a, *b) < 0 {
		return a
	} else {
		return b
	}
}

func (ws *Workspace) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	logGroups := input.Logs
	if logGroups == nil {
		// No filter specified, use all streams.
		logDescriptions, err := ws.DescribeLogs(ctx, &api.DescribeLogsInput{})
		if err != nil {
			return nil, fmt.Errorf("enumerating logs: %w", err)
		}
		logGroups = make([]string, len(logDescriptions.Logs))
		for i, group := range logDescriptions.Logs {
			logGroups[i] = group.Name
		}
	}
	logStreams := make([]string, 0, 2*len(logGroups))
	// Expand log groups in to streams.
	for _, group := range logGroups {
		// Each process acts as a log group combining both stdout and stderr.
		// XXX See note [LOG_COMPONENTS].
		for _, suffix := range []string{"", ":out", ":err"} {
			stream := group + suffix
			logStreams = append(logStreams, stream)
		}
	}

	collector := log.CurrentLogCollector(ctx)
	collectorOutput, err := collector.GetEvents(ctx, &logd.GetEventsInput{
		Logs:      logStreams,
		Cursor:    input.Cursor,
		FilterStr: input.FilterStr,
		Prev:      input.Prev,
		Next:      input.Next,
	})
	if err != nil {
		return nil, err
	}
	output := api.GetEventsOutput{
		Items:      make([]api.Event, len(collectorOutput.Items)),
		PrevCursor: collectorOutput.PrevCursor,
		NextCursor: collectorOutput.NextCursor,
	}
	for i, collectorEvent := range collectorOutput.Items {
		output.Items[i] = api.Event{
			ID:        collectorEvent.ID,
			Log:       collectorEvent.Log,
			Timestamp: collectorEvent.Timestamp,
			Message:   collectorEvent.Message,
		}
	}
	return &output, nil
}

func (ws *Workspace) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	jobID := ws.controlEachProcess(ctx, "starting", func(ctx context.Context, process api.Process) error {
		_, err := process.Start(ctx, &api.StartInput{})
		return err
	})
	return &api.StartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StartComponents(ctx context.Context, input *api.StartComponentsInput) (*api.StartComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "starting", filter, func(ctx context.Context, process api.Process) error {
		_, err := process.Start(ctx, &api.StartInput{})
		return err
	})
	return &api.StartComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	jobID := ws.controlEachProcess(ctx, "stopping", func(ctx context.Context, process api.Process) error {
		_, err := process.Stop(ctx, input)
		return err
	})
	return &api.StopOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StopComponents(ctx context.Context, input *api.StopComponentsInput) (*api.StopComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "stopping", filter, func(ctx context.Context, process api.Process) error {
		_, err := process.Stop(ctx, &api.StopInput{TimeoutSeconds: input.TimeoutSeconds})
		return err
	})
	return &api.StopComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Restart(ctx context.Context, input *api.RestartInput) (*api.RestartOutput, error) {
	jobID := ws.controlEachProcess(ctx, "restarting", func(ctx context.Context, process api.Process) error {
		_, err := process.Restart(ctx, input)
		return err
	})
	return &api.RestartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) RestartComponents(ctx context.Context, input *api.RestartComponentsInput) (*api.RestartComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "restart", filter, func(ctx context.Context, process api.Process) error {
		_, err := process.Restart(ctx, &api.RestartInput{TimeoutSeconds: input.TimeoutSeconds})
		return err
	})
	return &api.RestartComponentsOutput{
		JobID: jobID,
	}, nil
}

// TODO: Filter by interface, not concrete type.
var processTypes = []string{"process", "container"}

func (ws *Workspace) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (*api.DescribeProcessesOutput, error) {
	components, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{
		Types: processTypes,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	logger := logging.CurrentLogger(ctx)

	output := api.DescribeProcessesOutput{
		Processes: make([]api.ProcessDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		// XXX Violates component state encapsulation.
		switch component.Type {
		case "process":
			var state process.State
			if err := jsonutil.UnmarshalString(component.State, &state); err != nil {
				logger.Infof("unmarshalling process state: %v\n", err)
				continue
			}

			process := api.ProcessDescription{
				ID:       component.ID,
				Name:     component.Name,
				Provider: "unix",
				EnvVars:  state.FullEnvironment,
				Spec:     component.Spec,
			}

			proc, err := psprocess.NewProcess(int32(state.Pid))
			if err == nil {
				process.Running, err = proc.IsRunning()
				if err != nil {
					return nil, err
				}

				memoryInfo, err := proc.MemoryInfo()
				if err != nil {
					return nil, err
				}

				process.ResidentMemory = memoryInfo.RSS

				connections, err := proc.Connections()
				if err != nil {
					return nil, err
				}

				var ports []uint32
				for _, conn := range connections {
					if conn.Laddr.Port != 0 {
						ports = append(ports, conn.Laddr.Port)
					}
				}
				process.Ports = ports

				process.CreateTime, err = proc.CreateTime()
				if err != nil {
					return nil, err
				}

				children, err := proc.Children()
				if err == nil {
					var childrenExecutables []string
					for _, child := range children {
						exe, err := child.Exe()
						if err != nil {
							return nil, err
						}
						childrenExecutables = append(childrenExecutables, exe)
					}
					process.ChildrenExecutables = childrenExecutables
				}

				process.CPUPercent, err = proc.CPUPercent()
				if err != nil {
					return nil, err
				}
			}
			output.Processes = append(output.Processes, process)
		case "container":
			process, err := container.GetProcessDescription(ctx, ws.Docker, component)
			if err != nil {
				return nil, fmt.Errorf("could not get container process description: %w", err)
			}
			output.Processes = append(output.Processes, process)
		}
	}
	return &output, nil
}

func (ws *Workspace) DescribeVolumes(ctx context.Context, input *api.DescribeVolumesInput) (*api.DescribeVolumesOutput, error) {
	components, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{
		Types: []string{"volume"},
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	output := api.DescribeVolumesOutput{
		Volumes: make([]api.VolumeDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		volume := api.VolumeDescription{
			ID:   component.ID,
			Name: component.Name,
		}
		output.Volumes = append(output.Volumes, volume)
	}
	return &output, nil
}

func (ws *Workspace) DescribeNetworks(ctx context.Context, input *api.DescribeNetworksInput) (*api.DescribeNetworksOutput, error) {
	components, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{
		Types: []string{"network"},
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	output := api.DescribeNetworksOutput{
		Networks: make([]api.NetworkDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		network := api.NetworkDescription{
			ID:   component.ID,
			Name: component.Name,
		}
		output.Networks = append(output.Networks, network)
	}
	return &output, nil
}

func (ws *Workspace) ExportProcfile(ctx context.Context, input *api.ExportProcfileInput) (*api.ExportProcfileOutput, error) {
	logger := logging.CurrentLogger(ctx)
	procs, err := ws.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
	if err != nil {
		return nil, fmt.Errorf("describing processes: %w", err)
	}

	unixProcs := make([]procfile.Process, 0, len(procs.Processes))
	for _, proc := range procs.Processes {
		if proc.Provider == "unix" {
			var spec process.Spec
			if err := jsonutil.UnmarshalString(proc.Spec, &spec); err != nil {
				logger.Infof("unmarshalling process spec: %v\n", err)
				continue
			}

			unixProcs = append(unixProcs, procfile.Process{
				Name:        proc.Name,
				Program:     spec.Program,
				Arguments:   spec.Arguments,
				Environment: spec.Environment,
			})
		}
	}

	var export bytes.Buffer
	if err := procfile.Generate(&export, unixProcs); err != nil {
		return nil, fmt.Errorf("generating procfile: %w", err)
	}

	return &api.ExportProcfileOutput{
		Procfile: export.String(),
	}, nil
}

func (ws *Workspace) ReadFile(ctx context.Context, input *api.ReadFileInput) (*api.ReadFileOutput, error) {
	resolvedPath, err := ws.resolveWorkspacePath(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return &api.ReadFileOutput{
		Content: string(content),
	}, nil
}

func (ws *Workspace) WriteFile(ctx context.Context, input *api.WriteFileInput) (*api.WriteFileOutput, error) {
	resolvedPath, err := ws.resolveWorkspacePath(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	mode := os.FileMode(0644)
	if input.Mode != nil {
		mode = os.FileMode(*input.Mode)
	}

	if err := os.WriteFile(resolvedPath, []byte(input.Content), mode); err != nil {
		return nil, err
	}

	return &api.WriteFileOutput{}, nil
}

func (ws *Workspace) resolveWorkspacePath(ctx context.Context, relativePath string) (string, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return "", fmt.Errorf("describing workspace: %w", err)
	}

	resolvedPath := path.Join(description.Root, relativePath)
	if !pathutil.HasPathPrefix(resolvedPath, description.Root) {
		return "", fmt.Errorf("directory is outside workspace: %q", relativePath)
	}

	return resolvedPath, nil
}

type componentFilter struct {
	Refs  []string
	Types []string
}

func (ws *Workspace) controlEachComponent(ctx context.Context, label string, filter componentFilter, f interface{}) (jobID string) {
	job := ws.TaskTracker.StartTask(ctx, label)
	go func() {
		defer job.Finish()
		ws.goControlComponents(job, filter, f)
	}()
	return job.ID()
}

func (ws *Workspace) goControlComponents(t *task.Task, filter componentFilter, f interface{}) {
	describe := &api.DescribeComponentsInput{
		Types: filter.Types,
	}

	if filter.Refs != nil {
		ids, err := ws.resolveRefs(t, filter.Refs)
		if err != nil {
			t.Fail(fmt.Errorf("resolving refs: %w", err))
			return
		}
		describe.IDs = ids
	}

	components, err := ws.DescribeComponents(t, describe)
	if err != nil {
		t.Fail(fmt.Errorf("describing components: %w", err))
		return
	}

	for _, component := range components.Components {
		component := component
		t.Go(component.Name, func(t *task.Task) error {
			return ws.control(t, component, f)
		})
	}
}

func (ws *Workspace) controlEachProcess(ctx context.Context, label string, f interface{}) (jobID string) {
	filter := componentFilter{
		Types: processTypes,
	}
	return ws.controlEachComponent(ctx, label, filter, f)
}

func (ws *Workspace) control(ctx context.Context, desc api.ComponentDescription, f interface{}) error {
	ctrl := ws.newController(ctx, desc)
	if err := ctrl.InitResource(); err != nil {
		return err
	}
	ctxV := reflect.ValueOf(ctx)
	fV := reflect.ValueOf(f)
	ctrlV := reflect.ValueOf(ctrl)
	argT := fV.Type().In(1)
	if !ctrlV.Type().AssignableTo(argT) {
		return fmt.Errorf("%q controller does not implement %s", desc.Type, argT)
	}
	results := fV.Call([]reflect.Value{ctxV, ctrlV})
	fErr, _ := results[0].Interface().(error)
	// Try to save state even if f fails.
	newState, err := ctrl.MarshalState()
	if err == nil {
		if ctrl.IsDeleted() {
			// TODO: Do this asynchronously as a garbage collection pass, so
			// that we can inspect deleted components and debug them.
			_, err = ws.Store.RemoveComponent(ctx, &state.RemoveComponentInput{
				ID: desc.ID,
			})
		} else {
			_, err = ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
				ID:    desc.ID,
				State: newState,
			})
		}
	}
	if fErr != nil {
		return fErr
	}
	return err
}

func (ws *Workspace) Build(ctx context.Context, input *api.BuildInput) (*api.BuildOutput, error) {
	jobID := ws.controlEachProcess(ctx, "building", func(ctx context.Context, builder api.Builder) error {
		_, err := builder.Build(ctx, &api.BuildInput{})
		return err
	})
	return &api.BuildOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) BuildComponents(ctx context.Context, input *api.BuildComponentsInput) (*api.BuildComponentsOutput, error) {
	filter := componentFilter{
		Refs: input.Refs,
	}
	jobID := ws.controlEachComponent(ctx, "building", filter, func(ctx context.Context, builder api.Builder) error {
		_, err := builder.Build(ctx, &api.BuildInput{})
		return err
	})
	return &api.BuildComponentsOutput{
		JobID: jobID,
	}, nil
}
