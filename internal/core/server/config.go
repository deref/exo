package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/import/compose"
	"github.com/deref/exo/internal/import/procfile"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/pathutil"
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

func (ws *Workspace) loadManifest(rootDir string, input *api.ApplyInput) manifest.LoadResult {
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
					return manifest.LoadResult{
						Err: fmt.Errorf("searching for manifest: %w", err),
					}
				}
				if exist {
					manifestPath = candidatePath
					break
				}
			}
			if manifestPath == "" {
				return manifest.LoadResult{
					Err: errutil.NewHTTPError(http.StatusBadRequest, "could not find manifest file"),
				}
			}
		}

		if !pathutil.HasFilePathPrefix(manifestPath, rootDir) {
			return manifest.LoadResult{
				Err: errors.New("cannot read manifest outside of workspace root"),
			}
		}

		bs, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return manifest.LoadResult{
				Err: fmt.Errorf("reading manifest file: %w", err),
			}
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
			if strings.HasPrefix(name, "procfile.") || strings.HasSuffix(name, ".procfile") {
				format = "procfile"
			} else {
				return manifest.LoadResult{
					Err: errutil.NewHTTPError(http.StatusBadRequest, "cannot determine manifest format from file name"),
				}
			}
		}
	} else {
		format = *input.Format
	}

	var load func(r io.Reader) manifest.LoadResult
	switch format {
	case "procfile":
		load = procfile.Import
	case "compose":
		load = compose.Import
	case "exo":
		load = manifest.Read
	default:
		return manifest.LoadResult{
			Err: fmt.Errorf("unknown manifest format: %q", format),
		}
	}

	res := load(strings.NewReader(manifestString))
	if res.Err != nil {
		res.Err = errutil.WithHTTPStatus(http.StatusBadRequest, res.Err)
	}
	// TODO: Validate manifest.
	return res
}
