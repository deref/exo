package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/deref/exo/config"
	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/import/compose"
	"github.com/deref/exo/import/procfile"
	"github.com/deref/exo/util/osutil"
)

type configCandidate struct {
	Format   string
	Filename string
}

var configCandidates = []configCandidate{
	{"exo", "exo.hcl"},
	{"compose", "compose.yaml"},
	{"compose", "compose.yml"},
	{"compose", "docker-compose.yaml"},
	{"compose", "docker-compose.yml"},
	{"procfile", "Procfile"},
}

func (ws *Workspace) resolveConfig(rootDir string, input *api.ApplyInput) (*config.Config, error) {
	configString := ""
	configPath := ""
	if input.ConfigPath != nil {
		configPath = *input.ConfigPath
	}
	if input.Config == nil {
		if input.ConfigPath == nil {
			// Search for config.
			for _, candidate := range configCandidates {
				if input.Format != nil && *input.Format != candidate.Format {
					continue
				}
				candidatePath := filepath.Join(rootDir, candidate.Filename)
				exist, err := osutil.Exists(candidatePath)
				if err != nil {
					return nil, fmt.Errorf("searching for config: %w", err)
				}
				if exist {
					configPath = candidatePath
					break
				}
			}
			if configPath == "" {
				return nil, errors.New("could not find config file")
			}
		}

		if !filepath.HasPrefix(configPath, rootDir) {
			return nil, errors.New("cannot read config outside of workspace root")
		}

		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		configString = string(bs)
	} else {
		configString = *input.Config
	}

	format := ""
	if input.Format == nil {
		// Guess format.
		name := strings.ToLower(filepath.Base(configPath))
		switch name {
		case "procfile":
			format = "procfile"
		case "compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml":
			format = "compose"
		case "exo.hcl", "":
			format = "exo"
		default:
			if strings.HasSuffix(name, ".procfile") {
				format = "procfile"
			} else {
				return nil, errors.New("cannot determine config format from file name")
			}
		}
	} else {
		format = *input.Format
	}

	var load func(r io.Reader) (*config.Config, error)
	switch format {
	case "procfile":
		load = procfile.Import
	case "compose":
		load = compose.Import
	case "exo":
		load = config.Read
	default:
		return nil, fmt.Errorf("unknown config format: %q", format)
	}

	return load(strings.NewReader(configString))
}
