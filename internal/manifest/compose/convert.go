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
	res := manifest.LoadResult{
		Manifest: &manifest.Manifest{},
	}

	// TODO: Is there something like json.RawMessage so we can
	// avoid marshalling and re-marshalling each spec?
	for originalName, service := range comp.Services {
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
	for originalName, network := range comp.Networks {
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
	for originalName, volume := range comp.Volumes {
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
