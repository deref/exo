package statefile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/deref/exo/core/state/api"
	state "github.com/deref/exo/core/state/api"
	"github.com/deref/exo/util/atom"
	"github.com/deref/exo/util/errutil"
	"github.com/deref/exo/util/pathutil"
)

func New(filename string) *Store {
	return &Store{
		atom: atom.NewFileAtom(filename, atom.CodecJSON),
	}
}

type Store struct {
	atom atom.Atom
}

var _ state.Store = (*Store)(nil)

type Root struct {
	Workspaces          map[string]*Workspace `json:"workspaces"`          // Keyed by ID.
	ComponentWorkspaces map[string]string     `json:"componentWorkspaces"` // Component ID -> Workspace ID.
}

type Workspace struct {
	Root       string                `json:"root"`
	Names      map[string]string     `json:"names"`      // Name -> ID.
	Components map[string]*Component `json:"components"` // Keyed by ID.
}

type Component struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Spec        string  `json:"spec"`
	State       string  `json:"state"`
	Created     string  `json:"created"`
	Initialized *string `json:"initialized"`
	Disposed    *string `json:"disposed"`
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

func (sto *Store) FindWorkspace(ctx context.Context, input *state.FindWorkspaceInput) (*state.FindWorkspaceOutput, error) {
	var root Root
	if err := sto.atom.Deref(&root); err != nil {
		return nil, err
	}
	// Find the deepest root prefix match.
	maxLen := 0
	found := ""
	for id, workspace := range root.Workspaces {
		n := len(workspace.Root)
		if n > maxLen && pathutil.HasFilePathPrefix(input.Path, workspace.Root) {
			found = id
			maxLen = n
		}
	}
	var output state.FindWorkspaceOutput
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
	results := make([]*string, len(input.Refs))
	for i, ref := range input.Refs {
		if _, isID := workspace.Components[ref]; isID {
			id := ref
			results[i] = &id
			continue
		}
		id := workspace.Names[ref]
		if id != "" {
			results[i] = &id
		}
	}
	return &state.ResolveOutput{IDs: results}, nil
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
	}
	for _, id := range input.IDs {
		ids[id] = true
	}

	for componentID, component := range workspace.Components {
		if ids == nil || ids[componentID] {
			output.Components = append(output.Components, state.ComponentDescription{
				ID:          componentID,
				WorkspaceID: input.WorkspaceID,
				Name:        component.Name,
				Type:        component.Type,
				Spec:        component.Spec,
				State:       component.State,
				Created:     component.Created,
				Initialized: component.Initialized,
				Disposed:    component.Disposed,
			})
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
			Name:    input.Name,
			Type:    input.Type,
			Spec:    input.Spec,
			Created: input.Created,
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
