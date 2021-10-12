package procfile

import (
	"fmt"
	"io"
	"strconv"

	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/providers/unix/components/process"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/hashicorp/hcl/v2"
)

type loader struct{}

var Loader = loader{}

func (l loader) Load(r io.Reader) (*exohcl.Manifest, error) {
	procfile, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return convert(procfile)
}

const BasePort = 5000
const PortStep = 100

func convert(procfile *Procfile) (*exohcl.Manifest, error) {
	var diags hcl.Diagnostics
	m := &exohcl.Manifest{}
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
		name := exohcl.MangleName(p.Name)
		if name != p.Name {
			var subject *hcl.Range
			diags = append(diags, exohcl.NewRenameWarning(p.Name, name, subject))
		}

		// Add component.
		component := exohcl.Component{
			Name: name,
			Type: "process",
			Spec: jsonutil.MustMarshalIndentString(process.Spec{
				Program:     p.Program,
				Arguments:   p.Arguments,
				Environment: environment,
			}),
		}
		m.Components = append(m.Components, component)
	}
	return m, diags
}
