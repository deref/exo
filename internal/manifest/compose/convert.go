package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/yamlutil"
)

type Importer struct {
	// ProjectName is used as a prefix for the resources created by this importer.
	ProjectName string
}

func (i *Importer) Load(r io.Reader) manifest.LoadResult {
	composeProject, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return i.convert(composeProject)
}

func (i *Importer) convert(project *Project) manifest.LoadResult {
	res := manifest.LoadResult{
		Manifest: &manifest.Manifest{},
	}
	for originalName, service := range project.Services {
		// TODO: Create a higher level concept for a service. The service should be named
		// without a prefix, but the container should have the prefixed name.
		name, renamed := i.prefixedName(originalName)
		if renamed {
			res = res.AddRenameWarning(originalName, name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(service),
		})
	}
	hasDefaultNetwork := false
	for originalName, network := range project.Networks {
		if originalName == "default" {
			hasDefaultNetwork = true
		}
		name, renamed := i.prefixedName(originalName)
		if renamed {
			res = res.AddRenameWarning(originalName, name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(network),
		})
	}
	if !hasDefaultNetwork {
		name, _ := i.prefixedName("default")
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(map[string]string{
				"driver": "bridge",
			}),
		})
	}
	for originalName, volume := range project.Volumes {
		name, renamed := i.prefixedName(originalName)
		if renamed {
			res = res.AddRenameWarning(originalName, name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "volume",
			Spec: yamlutil.MustMarshalString(volume),
		})
	}
	return res
}

func (i *Importer) prefixedName(name string) (string, bool) {
	newName := manifest.MangleName(name)
	renamed := newName != name
	return fmt.Sprintf("%s_%s", i.ProjectName, newName), renamed
}
