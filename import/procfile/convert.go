package procfile

import (
	"fmt"
	"io"
	"strconv"

	"github.com/deref/exo/manifest"
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
	for _, process := range procfile.Processes {
		component := manifest.Component{
			Name: process.Name,
			Type: "process",
			Spec: jsonutil.MustMarshalString(map[string]interface{}{
				"program":   process.Program,
				"arguments": process.Arguments,
				"environment": map[string]interface{}{
					"PORT": strconv.Itoa(port),
				},
			}),
		}
		port += PortStep
		m.Components = append(m.Components, component)
	}
	return m, nil
}
