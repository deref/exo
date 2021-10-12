package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/procfile"
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
	{"compose", "compose.yaml"},
	{"compose", "compose.yml"},
	{"compose", "docker-compose.yaml"},
	{"compose", "docker-compose.yml"},
	{"procfile", "Procfile"},
}

func (ws *Workspace) loadManifest(rootDir string, input *api.ApplyInput) (*manifest.Manifest, error) {
	manifestString := ""
	manifestPath := ""
	if input.ManifestPath != nil {
		manifestPath = *input.ManifestPath
	}
	if input.Manifest == nil {
		if input.ManifestPath == nil {
			var err error
			manifestPath, err = ws.resolveManifest(rootDir, input.Format)
			if err != nil {
				return nil, err
			}
			if manifestPath == "" {
				return nil, errutil.NewHTTPError(http.StatusBadRequest, "could not find manifest file")
			}
		}

		if !pathutil.HasFilePathPrefix(manifestPath, rootDir) {
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
	if input.Format == "" {
		if manifestPath == "" {
			format = "exo"
		} else {
			format = manifest.GuessFormat(manifestPath)
			if format == "" {
				return nil, errutil.NewHTTPError(http.StatusBadRequest, "cannot determine manifest format from file name")
			}
		}
	} else {
		format = input.Format
	}

	var loader interface {
		Load(r io.Reader) (*manifest.Manifest, error)
	}
	switch format {
	case "procfile":
		loader = procfile.Loader
	case "compose":
		projectName := path.Base(rootDir)
		projectName = manifest.MangleName(projectName)
		loader = &compose.Loader{ProjectName: projectName}
	case "exo":
		loader = &manifest.Loader{
			Filename: manifestPath,
		}
	default:
		return nil, fmt.Errorf("unknown manifest format: %q", format)
	}

	return loader.Load(strings.NewReader(manifestString))
}

func (ws *Workspace) ResolveManifest(ctx context.Context, input *api.ResolveManifestInput) (*api.ResolveManifestOutput, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return nil, fmt.Errorf("describing workspace: %w", err)
	}
	path, err := ws.resolveManifest(description.Root, input.Format)
	if err != nil {
		return nil, err
	}
	return &api.ResolveManifestOutput{
		Path: path,
	}, nil
}

func (ws *Workspace) resolveManifest(rootDir, format string) (string, error) {
	for _, candidate := range manifestCandidates {
		if format != "" && format != candidate.Format {
			continue
		}
		candidatePath := filepath.Join(rootDir, candidate.Filename)
		exist, err := osutil.Exists(candidatePath)
		if err != nil {
			return "", fmt.Errorf("searching for manifest: %w", err)
		}
		if exist {
			return candidatePath, nil
		}
	}
	return "", nil
}
