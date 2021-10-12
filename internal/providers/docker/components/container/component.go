package container

import (
	"path"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/docker/docker/api/types/strslice"
)

type Container struct {
	docker.ComponentBase

	State State

	SyslogPort uint
}

func (c *Container) ProjectName() string {
	projectName := path.Base(c.WorkspaceRoot)
	projectName = exohcl.MangleName(projectName)
	return projectName
}

type Spec compose.Service

type State struct {
	ContainerID string     `json:"containerId"`
	Running     bool       `json:"running"`
	Image       ImageState `json:"image"`
}

type ImageState struct {
	ID         string            `json:"id"`
	Spec       string            `json:"spec"`
	Command    strslice.StrSlice `json:"command"`
	WorkingDir string            `json:"workingDir"`
	Entrypoint strslice.StrSlice `json:"entrypoint"`
	Shell      []string          `json:"shell"`
}
