package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/yamlutil"
)

func Import(r io.Reader) manifest.LoadResult {
	procfile, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return Convert(procfile)
}

func Convert(project *Project) manifest.LoadResult {
	res := manifest.LoadResult{
		Manifest: &manifest.Manifest{},
	}
	for originalName, service := range project.Services {
		name := manifest.MangleName(originalName)
		if name != originalName {
			res = res.AddRenameWarning(originalName, name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(service),
		})
	}
	for originalName, network := range project.Networks {
		name := manifest.MangleName(originalName)
		if name != originalName {
			res = res.AddRenameWarning(originalName, name)
		}
		res.Manifest.Components = append(res.Manifest.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(network),
		})
	}
	for originalName, volume := range project.Volumes {
		name := manifest.MangleName(originalName)
		if name != originalName {
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
