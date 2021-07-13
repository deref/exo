package statefile

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/atom"
	"github.com/deref/exo/state"
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
	Projects          map[string]*Project `json:"projects"`          // Keyed by ID.
	ComponentProjects map[string]string   `json:"componentProjects"` // Component ID -> Project ID.
}

type Project struct {
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
		if root.Projects == nil {
			root.Projects = make(map[string]*Project)
		}
		if root.ComponentProjects == nil {
			root.ComponentProjects = make(map[string]string)
		}
		return f(&root)
	})
	return &root, err
}

func (sto *Store) DescribeComponents(ctx context.Context, input *state.DescribeComponentsInput) (*state.DescribeComponentsOutput, error) {
	if input.ProjectID == "" {
		return nil, errors.New("project-id is required")
	}

	root, err := sto.deref()
	if err != nil {
		return nil, err
	}

	output := &state.DescribeComponentsOutput{
		Components: []state.ComponentDescription{},
	}

	var project *Project
	if root.Projects != nil {
		project = root.Projects[input.ProjectID]
	}
	if project == nil {
		return output, nil
	}

	var names map[string]bool
	if input.Names != nil {
		names = make(map[string]bool, len(input.Names))
	}
	for _, name := range input.Names {
		names[name] = true
	}

	for componentID, component := range project.Components {
		if names == nil || names[component.Name] {
			output.Components = append(output.Components, state.ComponentDescription{
				ID:          componentID,
				ProjectID:   input.ProjectID,
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
	return output, nil
}

func (sto *Store) AddComponent(ctx context.Context, input *state.AddComponentInput) (*state.AddComponentOutput, error) {
	if input.ProjectID == "" {
		return nil, errors.New("project-id is required")
	}
	_, err := sto.swap(func(root *Root) error {
		project := root.Projects[input.ProjectID]
		if project == nil {
			project = &Project{}
			root.Projects[input.ProjectID] = project
		}
		if project.Components == nil {
			project.Components = make(map[string]*Component)
		}
		if project.Names == nil {
			project.Names = make(map[string]string)
		}
		if project.Names[input.Name] != "" {
			return fmt.Errorf("component named %q already exits", input.Name)
		}
		project.Names[input.Name] = input.ID
		if project.Components[input.ID] != nil {
			return fmt.Errorf("component id %q already exists", input.ID)
		}
		project.Components[input.ID] = &Component{
			Name:    input.Name,
			Type:    input.Type,
			Spec:    input.Spec,
			Created: input.Created,
		}
		root.ComponentProjects[input.ID] = input.ProjectID
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
		projectId := root.ComponentProjects[input.ID]
		if projectId == "" {
			return errors.New("cannot find project for component")
		}
		project := root.Projects[projectId]
		if project == nil {
			return errors.New("corrupt state: no project for component")
		}
		component := project.Components[input.ID]
		if component == nil {
			return errors.New("corrupt state: component not in project")
		}
		if input.Disposed != "" {
			// TODO: Validate disposed is date.
			component.Disposed = &input.Disposed
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
		projectID := root.ComponentProjects[input.ID]
		if projectID == "" {
			return fmt.Errorf("cannot find project for component %q", input.ID)
		}
		project := root.Projects[projectID]
		if project == nil {
			return fmt.Errorf("component %q has invalid project %q", input.ID, projectID)
		}
		var component *Component
		if project.Components != nil {
			component = project.Components[input.ID]
		}
		if component == nil {
			return fmt.Errorf("no component for id: %q", input.ID)
		}
		delete(project.Components, input.ID)
		if project.Names != nil {
			delete(project.Names, component.Name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state.RemoveComponentOutput{}, nil
}
