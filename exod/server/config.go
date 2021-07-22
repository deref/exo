package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/import/compose"
	"github.com/deref/exo/import/procfile"
	"github.com/deref/exo/manifest"
	"github.com/deref/exo/util/errutil"
	"github.com/deref/exo/util/osutil"
)

type manifestCandidate struct {
	Format   string
	Filename string
}

var manifestCandidates = []manifestCandidate{
	{"exo", "exo.hcl"},
	// TODO: Uncomment when we have docker-compose support.
	//{"compose", "compose.yaml"},
	//{"compose", "compose.yml"},
	//{"compose", "docker-compose.yaml"},
	//{"compose", "docker-compose.yml"},
	{"procfile", "Procfile"},
}

func (ws *Workspace) resolveManifest(rootDir string, input *api.ApplyInput) (*manifest.Manifest, error) {
	manifestString := ""
	manifestPath := ""
	if input.ManifestPath != nil {
		manifestPath = *input.ManifestPath
	}
	if input.Manifest == nil {
		if input.ManifestPath == nil {
			// Search for manifest.
			for _, candidate := range manifestCandidates {
				if input.Format != nil && *input.Format != candidate.Format {
					continue
				}
				candidatePath := filepath.Join(rootDir, candidate.Filename)
				exist, err := osutil.Exists(candidatePath)
				if err != nil {
					return nil, fmt.Errorf("searching for manifest: %w", err)
				}
				if exist {
					manifestPath = candidatePath
					break
				}
			}
			if manifestPath == "" {
				return nil, errutil.NewHTTPError(http.StatusBadRequest, "could not find manifest file")
			}
		}

		if !filepath.HasPrefix(manifestPath, rootDir) {
			return nil, errors.New("cannot read manifest outside of workspace root")
		}

		bs, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("reading manifest file: %w", err)
		}
		manifestString = string(bs)
	} else {
		manifestString = *input.Manifest
	}

	format := ""
	if input.Format == nil {
		// Guess format.
		name := strings.ToLower(filepath.Base(manifestPath))
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
				return nil, errors.New("cannot determine manifest format from file name")
			}
		}
	} else {
		format = *input.Format
	}

	var load func(r io.Reader) (*manifest.Manifest, error)
	switch format {
	case "procfile":
		load = procfile.Import
	case "compose":
		load = compose.Import
	case "exo":
		load = manifest.Read
	default:
		return nil, fmt.Errorf("unknown manifest format: %q", format)
	}

	return load(strings.NewReader(manifestString))
}
