package procfile

import (
	"fmt"
	"io"
	"strconv"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/providers/unix/components/process"
	"github.com/deref/exo/internal/util/jsonutil"
)

func Import(r io.Reader) manifest.LoadResult {
	procfile, err := Parse(r)
	if err != nil {
		return manifest.LoadResult{Err: fmt.Errorf("parsing: %w", err)}
	}
	return Convert(procfile)
}

const BasePort = 5000
const PortStep = 100

func Convert(procfile *Procfile) manifest.LoadResult {
	var res manifest.LoadResult
	m := manifest.NewManifest()
	port := BasePort
	for _, p := range procfile.Processes {
		// Assign default PORT, merge in specified environment.
		environment := map[string]string{
			"PORT": strconv.Itoa(port),
		}
		for name, value := range p.Environment {
			environment[name] = value
		}
		port += PortStep

		// Get component name.
		name := manifest.MangleName(p.Name)
		if name != p.Name {
			warning := fmt.Sprintf("invalid name: %q, renamed to: %q", p.Name, name)
			res.Warnings = append(res.Warnings, warning)
		}

		// Add component.
		component := manifest.Component{
			Name: name,
			Type: "process",
			Spec: jsonutil.MustMarshalString(process.Spec{
				Program:     p.Program,
				Arguments:   p.Arguments,
				Environment: environment,
			}),
		}
		m.Components = append(m.Components, component)
	}
	res.Manifest = m
	return res
}
