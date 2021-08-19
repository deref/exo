package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/yamlutil"
)

type Loader struct {
	// ProjectName is used as a prefix for the resources created by this importer.
	ProjectName string
}

func (i *Loader) Load(r io.Reader) manifest.LoadResult {
	composeProject, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return i.convert(composeProject)
}

func (i *Loader) convert(project *Project) manifest.LoadResult {
	res := manifest.LoadResult{
		Manifest: &manifest.Manifest{},
	}
	for originalName, service := range project.Services {
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}
		s := service.ToMap()
		if _, ok := s["container_name"]; !ok {
			s["container_name"] = i.prefixedName(name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(s),
		})
	}
	hasDefaultNetwork := false
	for originalName, network := range project.Networks {
		if originalName == "default" {
			hasDefaultNetwork = true
		}
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}
		n := network.ToMap()

		// If a `name` key is specified in the network configuration (usually used in conjunction with `external: true`),
		// then we should honor that as the docker network name. Otherwise, we should set the name as
		// `<project_name>_<network_key>`.
		if _, ok := n["name"]; !ok {
			n["name"] = i.prefixedName(name)
		}

		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(n),
		})
	}
	if !hasDefaultNetwork {
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: "default",
			Type: "network",
			Spec: yamlutil.MustMarshalString(map[string]string{
				"name":   i.prefixedName("default"),
				"driver": "bridge",
			}),
		})
	}
	for originalName, volume := range project.Volumes {
		name := manifest.MangleName(originalName)
		if originalName != name {
			res = res.AddRenameWarning(originalName, name)
		}
		v := volume.ToMap()
		if _, ok := v["name"]; !ok {
			v["name"] = i.prefixedName(name)
		}
		// TODO
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "volume",
			Spec: yamlutil.MustMarshalString(v),
		})
	}
	return res
}

func (i *Loader) prefixedName(name string) string {
	return fmt.Sprintf("%s_%s", i.ProjectName, name)
}
