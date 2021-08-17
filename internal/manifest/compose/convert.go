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
	var m manifest.Manifest
	for name, service := range project.Services {
		m.Components = append(m.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(service),
		})
	}
	for name, network := range project.Networks {
		m.Components = append(m.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(network),
		})
	}
	for name, volume := range project.Volumes {
		m.Components = append(m.Components, manifest.Component{
			Name: name,
			Type: "volume",
			Spec: yamlutil.MustMarshalString(volume),
		})
	}
	return manifest.LoadResult{
		Manifest: &m,
	}
}
