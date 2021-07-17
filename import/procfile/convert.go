package procfile

import (
	"fmt"
	"io"
	"strconv"

	"github.com/deref/exo/config"
	"github.com/deref/exo/jsonutil"
)

func Import(r io.Reader) (*config.Config, error) {
	procfile, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return Convert(procfile)
}

const BasePort = 5000
const PortStep = 100

func Convert(procfile *Procfile) (*config.Config, error) {
	cfg := config.NewConfig()
	port := BasePort
	for _, process := range procfile.Processes {
		component := config.Component{
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
		cfg.Components = append(cfg.Components, component)
	}
	return cfg, nil
}
