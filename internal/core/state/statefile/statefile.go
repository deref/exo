package statefile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/deref/exo/internal/core/state/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/deps"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/atom"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/pathutil"
)

type Config struct {
	StoreFilename    string
	DeviceIDFilename string
}

func New(cfg Config) *Store {
	return &Store{
		atom:             atom.NewFileAtom(cfg.StoreFilename, atom.CodecJSON),
		deviceIDFilename: cfg.DeviceIDFilename,
	}
}

type Store struct {
	atom             atom.Atom
	deviceIDFilename string
}

var _ state.Store = (*Store)(nil)

type Root struct {
	Workspaces          map[string]*Workspace `json:"workspaces"`          // Keyed by ID.
	ComponentWorkspaces map[string]string     `json:"componentWorkspaces"` // Component ID -> Workspace ID.
	DeviceID            string                `json:"deviceId"`
}

type Component struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Spec        string   `json:"spec"`
	State       string   `json:"state"`
	Created     string   `json:"created"`
	Initialized *string  `json:"initialized"`
	Disposed    *string  `json:"disposed"`
	DependsOn   []string `json:"dependsOn"`
}

func (c *Component) getDescription(id, workspaceID string) state.ComponentDescription {
	return state.ComponentDescription{
		ID:          id,
		WorkspaceID: workspaceID,
		Name:        c.Name,
		Type:        c.Type,
		Spec:        c.Spec,
		State:       c.State,
		Created:     c.Created,
		Initialized: c.Initialized,
		Disposed:    c.Disposed,
		DependsOn:   c.DependsOn,
	}
}

func (sto *Store) deref() (*Root, error) {
	var root Root
	err := sto.atom.Deref(&root)
	return &root, err
}

func (sto *Store) swap(f func(root *Root) error) (*Root, error) {
	var root Root
	err := sto.atom.Swap(&root, func() error {
		if root.Workspaces == nil {
			root.Workspaces = make(map[string]*Workspace)
		}
		if root.ComponentWorkspaces == nil {
			root.ComponentWorkspaces = make(map[string]string)
		}
		return f(&root)
	})
	return &root, err
}

func (sto *Store) DescribeWorkspaces(ctx context.Context, input *state.DescribeWorkspacesInput) (*state.DescribeWorkspacesOutput, error) {
	var root Root
	if err := sto.atom.Deref(&root); err != nil {
		return nil, err
	}

	var ids map[string]bool
	if input.IDs != nil {
		ids = make(map[string]bool, len(input.IDs))
	}
	for _, id := range input.IDs {
		ids[id] = true
	}

	var output state.DescribeWorkspacesOutput
	for id, workspace := range root.Workspaces {
		if ids == nil || ids[id] {
			output.Workspaces = append(output.Workspaces, state.WorkspaceDescription{
				ID:   id,
				Root: workspace.Root,
			})
		}
	}
	return &output, nil
}

func (sto *Store) AddWorkspace(ctx context.Context, input *state.AddWorkspaceInput) (*state.AddWorkspaceOutput, error) {
	rootPath := filepath.Clean(input.Root)
	if !filepath.IsAbs(rootPath) {
		return nil, errutil.NewHTTPError(http.StatusBadRequest, "root must be absolute path")
	}
	_, err := sto.swap(func(root *Root) error {
		if root.Workspaces == nil {
			root.Workspaces = make(map[string]*Workspace)
		}
		if _, exists := root.Workspaces[input.ID]; exists {
			return fmt.Errorf("workspace %q already exists", input.ID)
		}
		for _, workspace := range root.Workspaces {
			if rootPath == workspace.Root {
				return errutil.HTTPErrorf(http.StatusConflict, "workspace with root %q already exists", rootPath)
			}
		}
		root.Workspaces[input.ID] = &Workspace{
			Root: rootPath,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.AddWorkspaceOutput{}, nil
}

func (sto *Store) RemoveWorkspace(ctx context.Context, input *state.RemoveWorkspaceInput) (*state.RemoveWorkspaceOutput, error) {
	_, err := sto.swap(func(root *Root) error {
		if root.Workspaces == nil {
			root.Workspaces = make(map[string]*Workspace)
		}
		workspace := root.Workspaces[input.ID]
		if workspace == nil {
			return nil
		}
		if len(workspace.Components) > 0 {
			return fmt.Errorf("cannot remove non-empty workspace %q", input.ID)
		}
		delete(root.Workspaces, input.ID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.RemoveWorkspaceOutput{}, nil
}

func (sto *Store) ResolveWorkspace(ctx context.Context, input *state.ResolveWorkspaceInput) (*state.ResolveWorkspaceOutput, error) {
	var root Root
	if err := sto.atom.Deref(&root); err != nil {
		return nil, err
	}

	// Resolve by ID.
	if _, validID := root.Workspaces[input.Ref]; validID {
		return &state.ResolveWorkspaceOutput{
			ID: &input.Ref,
		}, nil
	}

	// Resolve by path. Searches for the deepest root prefix match.
	maxLen := 0
	found := ""
	for id, workspace := range root.Workspaces {
		n := len(workspace.Root)
		if n > maxLen && pathutil.HasFilePathPrefix(input.Ref, workspace.Root) {
			found = id
			maxLen = n
		}
	}
	var output state.ResolveWorkspaceOutput
	if maxLen > 0 {
		output.ID = &found
	}
	return &output, nil
}

func (sto *Store) Resolve(ctx context.Context, input *state.ResolveInput) (*state.ResolveOutput, error) {
	if input.WorkspaceID == "" {
		return nil, errors.New("workspace-id is required")
	}
	var root Root
	if err := sto.atom.Deref(&root); err != nil {
		return nil, err
	}
	workspace := root.Workspaces[input.WorkspaceID]
	if workspace == nil {
		return nil, fmt.Errorf("no such workspace: %q", input.WorkspaceID)
	}
	ids := workspace.resolve(input.Refs)

	return &state.ResolveOutput{IDs: ids}, nil
}

func (sto *Store) DescribeComponents(ctx context.Context, input *state.DescribeComponentsInput) (*state.DescribeComponentsOutput, error) {
	if input.WorkspaceID == "" {
		return nil, errors.New("workspace-id is required")
	}

	root, err := sto.deref()
	if err != nil {
		return nil, err
	}

	output := &state.DescribeComponentsOutput{
		Components: []state.ComponentDescription{},
	}

	var workspace *Workspace
	if root.Workspaces == nil {
		return nil, fmt.Errorf("no such workspace: %q", input.WorkspaceID)
	}
	workspace = root.Workspaces[input.WorkspaceID]
	if workspace == nil {
		return output, nil
	}

	var ids map[string]bool
	if input.IDs != nil {
		ids = make(map[string]bool, len(input.IDs))
		for _, id := range input.IDs {
			ids[id] = true
		}
	}

	var types map[string]bool
	if input.Types != nil {
		types = make(map[string]bool, len(input.Types))
		for _, typ := range input.Types {
			types[typ] = true
		}
	}

	var componentGraph *deps.Graph
	if input.IncludeDependencies || input.IncludeDependents {
		componentGraph = deps.New()
	}

	for componentID, component := range workspace.Components {
		if (ids == nil || ids[componentID]) &&
			(types == nil || types[component.Type]) {
			output.Components = append(output.Components, component.getDescription(componentID, input.WorkspaceID))
		}

		if input.IncludeDependencies || input.IncludeDependents {
			// Add component dependencies to graph.
			dependencyIDs := workspace.resolve(component.DependsOn)
			for _, dependencyID := range dependencyIDs {
				if dependencyID != nil {
					componentGraph.DependOn(deps.StringNode(componentID), deps.StringNode(*dependencyID))
				}
			}
		}
	}

	if input.IncludeDependencies || input.IncludeDependents {
		seen := make(map[string]struct{}, len(output.Components))
		nextIDs := []string{}
		markSeen := func(id string) {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				nextIDs = append(nextIDs, id)
			}
		}

		// Start search with all components that would be returned anyway.
		for _, component := range output.Components {
			markSeen(component.ID)
		}

		for len(nextIDs) > 0 {
			// Swap the list for the next iteration with a new list.
			ids := nextIDs
			nextIDs = []string{}

			// Discover more dependencies/dependents.
			if input.IncludeDependencies {
				for _, id := range ids {
					dependencies := componentGraph.Dependencies(id)
					for dependencyID := range dependencies {
						markSeen(dependencyID)
					}
				}
			}
			if input.IncludeDependents {
				for _, id := range ids {
					dependents := componentGraph.Dependents(id)
					for dependencyID := range dependents {
						markSeen(dependencyID)
					}
				}
			}
		}

		// Remove the components that we have already added to the output from the seen list.
		for _, component := range output.Components {
			delete(seen, component.ID)
		}

		// Resolve the remaining ids to components.
		for resolvedID := range seen {
			component := workspace.Components[resolvedID]
			output.Components = append(output.Components, component.getDescription(resolvedID, input.WorkspaceID))
		}
	}

	sort.Sort(componentsSort{output.Components})
	return output, nil
}

type componentsSort struct {
	components []api.ComponentDescription
}

func (iface componentsSort) Len() int {
	return len(iface.components)
}

func (iface componentsSort) Less(i, j int) bool {
	return strings.Compare(iface.components[i].Name, iface.components[j].Name) < 0
}

func (iface componentsSort) Swap(i, j int) {
	tmp := iface.components[i]
	iface.components[i] = iface.components[j]
	iface.components[j] = tmp
}

func (sto *Store) AddComponent(ctx context.Context, input *state.AddComponentInput) (*state.AddComponentOutput, error) {
	if input.WorkspaceID == "" {
		return nil, errors.New("workspace-id is required")
	}
	_, err := sto.swap(func(root *Root) error {
		workspace := root.Workspaces[input.WorkspaceID]
		if workspace == nil {
			return fmt.Errorf("no such workspace: %q", input.WorkspaceID)
		}
		if workspace.Components == nil {
			workspace.Components = make(map[string]*Component)
		}
		if workspace.Names == nil {
			workspace.Names = make(map[string]string)
		}
		if workspace.Names[input.Name] != "" {
			return errutil.HTTPErrorf(http.StatusConflict, "component named %q already exists", input.Name)
		}
		workspace.Names[input.Name] = input.ID
		if workspace.Components[input.ID] != nil {
			return fmt.Errorf("component id %q already exists", input.ID)
		}
		workspace.Components[input.ID] = &Component{
			Name:      input.Name,
			Type:      input.Type,
			Spec:      input.Spec,
			Created:   input.Created,
			DependsOn: input.DependsOn,
		}
		root.ComponentWorkspaces[input.ID] = input.WorkspaceID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.AddComponentOutput{}, nil
}

func (sto *Store) PatchComponent(ctx context.Context, input *state.PatchComponentInput) (*state.PatchComponentOutput, error) {
	if input.ID == "" {
		return nil, errors.New("component id is required")
	}
	_, err := sto.swap(func(root *Root) error {
		workspaceId := root.ComponentWorkspaces[input.ID]
		if workspaceId == "" {
			return errors.New("cannot find workspace for component")
		}
		workspace := root.Workspaces[workspaceId]
		if workspace == nil {
			return errors.New("corrupt state: no workspace for component")
		}
		component := workspace.Components[input.ID]
		if component == nil {
			return errors.New("corrupt state: component not in workspace")
		}
		if input.Initialized != "" {
			component.Initialized = &input.Initialized // TODO: Validate iso8601.
		}
		if input.Disposed != "" {
			component.Disposed = &input.Disposed // TODO: Validate iso8601.
		}
		if input.DependsOn != nil {
			component.DependsOn = *input.DependsOn
		}
		if input.State != "" {
			component.State = input.State
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.PatchComponentOutput{}, nil

}

func (sto *Store) RemoveComponent(ctx context.Context, input *state.RemoveComponentInput) (*state.RemoveComponentOutput, error) {
	_, err := sto.swap(func(root *Root) error {
		workspaceID := root.ComponentWorkspaces[input.ID]
		if workspaceID == "" {
			return fmt.Errorf("cannot find workspace for component %q", input.ID)
		}
		workspace := root.Workspaces[workspaceID]
		if workspace == nil {
			return fmt.Errorf("component %q has invalid workspace %q", input.ID, workspaceID)
		}
		var component *Component
		if workspace.Components != nil {
			component = workspace.Components[input.ID]
		}
		if component == nil {
			return fmt.Errorf("no component for id: %q", input.ID)
		}
		delete(workspace.Components, input.ID)
		if workspace.Names != nil {
			delete(workspace.Names, component.Name)
		}
		if root.ComponentWorkspaces != nil {
			delete(root.ComponentWorkspaces, input.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.RemoveComponentOutput{}, nil
}

func (sto *Store) EnsureDevice(ctx context.Context, input *state.EnsureDeviceInput) (*state.EnsureDeviceOutput, error) {
	var deviceID string

	_, err := sto.swap(func(root *Root) error {
		if root.DeviceID != "" {
			deviceID = root.DeviceID
			return nil
		}
		deviceIDFile, err := os.ReadFile(sto.deviceIDFilename)
		switch {
		case os.IsNotExist(err):
			deviceID = gensym.RandomBase32()
		case err != nil:
			return err
		default:
			deviceID = string(deviceIDFile)
		}

		root.DeviceID = deviceID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &state.EnsureDeviceOutput{
		DeviceID: deviceID,
	}, nil
}
