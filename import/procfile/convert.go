package procfile

import (
	"fmt"
	"io"
	"strconv"

	"github.com/deref/exo/manifest"
	"github.com/deref/exo/providers/unix/components/process"
	"github.com/deref/exo/util/jsonutil"
)

func Import(r io.Reader) (*manifest.Manifest, error) {
	procfile, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return Convert(procfile)
}

const BasePort = 5000
const PortStep = 100

func Convert(procfile *Procfile) (*manifest.Manifest, error) {
	m := manifest.NewManifest()
	port := BasePort
	for _, p := range procfile.Processes {
		component := manifest.Component{
			Name: p.Name,
			Type: "process",
			Spec: jsonutil.MustMarshalString(process.Spec{
				Program:   p.Program,
				Arguments: p.Arguments,
				Environment: map[string]string{
					"PORT": strconv.Itoa(port),
				},
			}),
		}
		port += PortStep
		m.Components = append(m.Components, component)
	}
	return m, nil
}
