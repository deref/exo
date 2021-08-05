package compose

import (
	"fmt"
	"io"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/yamlutil"
)

func Import(r io.Reader) manifest.LoadResult {
	procfile, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return Convert(procfile)
}

func Convert(comp *compose.Compose) manifest.LoadResult {
	var m manifest.Manifest
	// TODO: Is there something like json.RawMessage so we can
	// avoid marshalling and re-marshalling each spec?
	for name, service := range comp.Services {
		m.Components = append(m.Components, manifest.Component{
			Name: name,
			Type: "container",
			Spec: yamlutil.MustMarshalString(service),
		})
	}
	for name, network := range comp.Networks {
		m.Components = append(m.Components, manifest.Component{
			Name: name,
			Type: "network",
			Spec: yamlutil.MustMarshalString(network),
		})
	}
	for name, volume := range comp.Volumes {
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
