package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/pathutil"
	"github.com/natefinch/atomic"
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

// XXX This is a hacky workaround until manifest handling is overhauled.
func (ws *Workspace) tryLoadManifest(ctx context.Context) *exohcl.Manifest {
	wsDesc, err := ws.describe(ctx)
	if err != nil {
		return nil
	}
	m, _ := ws.loadManifest(ctx, wsDesc.Root, &api.ApplyInput{})
	return m
}

func (ws *Workspace) loadManifest(ctx context.Context, rootDir string, input *api.ApplyInput) (*exohcl.Manifest, error) {
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

	// TODO: Get official name from workspace description.
	workspaceName := path.Base(rootDir)
	workspaceName = exohcl.MangleName(workspaceName)

	analysisContext := &exohcl.AnalysisContext{
		Context: ctx,
	}
	loader := &manifest.Loader{
		WorkspaceName: workspaceName,
		Format:        input.Format,
		Filename:      manifestPath,
		Bytes:         []byte(manifestString),
	}
	m, err := loader.Load(analysisContext)
	if err == nil && len(analysisContext.Diagnostics) > 0 {
		err = analysisContext.Diagnostics
	}
	return m, err
}

func (ws *Workspace) modifyManifest(ctx context.Context, re exohcl.Rewrite) error {
	wsDesc, err := ws.describe(ctx)
	if err != nil {
		return fmt.Errorf("describing workspace: %w", err)
	}
	manifest, err := ws.loadManifest(ctx, wsDesc.Root, &api.ApplyInput{})
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}
	f := exohcl.RewriteManifest(re, manifest)
	var buf bytes.Buffer
	if _, err := hclgen.WriteTo(&buf, f); err != nil {
		return fmt.Errorf("generating new manifest: %w", err)
	}
	manifestPath := filepath.Join(wsDesc.Root, "exo.hcl")
	if err := atomic.WriteFile(manifestPath, &buf); err != nil {
		return fmt.Errorf("writing new manifest: %w", err)
	}
	return nil
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
