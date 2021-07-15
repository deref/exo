package procfile

import (
	"fmt"
	"io"

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

func Convert(procfile *Procfile) (*config.Config, error) {
	cfg := config.NewConfig()
	for _, process := range procfile.Processes {
		component := config.Component{
			Name: process.Name,
			Type: "process",
			Spec: jsonutil.MustMarshalString(map[string]interface{}{
				"command":   process.Command,
				"arguments": process.Arguments,
			}),
		}
		cfg.Components = append(cfg.Components, component)
	}
	return cfg, nil
}
