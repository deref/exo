package resolvers

import (
	"context"
	"fmt"
	"strings"

	"github.com/deref/exo/internal/util/cueutil"
	"github.com/natefinch/atomic"
)

type ManifestResolver struct {
	Format string
	File   *FileResolver
	// Used if File is nil for in-memory manifests.
	content string
}

func (r *QueryResolver) MakeManifest(ctx context.Context, args struct {
	Content string
	Format  *string
}) *ManifestResolver {
	manifest := &ManifestResolver{
		content: args.Content,
	}
	if args.Format == nil {
		manifest.Format = defaultManifestFormat
	} else {
		manifest.Format = *args.Format
	}
	return manifest
}

func (r *QueryResolver) manifestByPath(ctx context.Context, fs *FileSystemResolver, path string, format *string) (*ManifestResolver, error) {
	file, err := fs.file(ctx, path)
	if err != nil {
		return nil, err
	}
	manifest := &ManifestResolver{
		File: file,
	}
	if format == nil {
		manifest.Format = defaultManifestFormat
	} else {
		manifest.Format = *format
	}
	return manifest, nil
}

const defaultManifestFormat = "exo"

type manifestCandidate struct {
	Format   string
	Filename string
}

var manifestCandidates = []manifestCandidate{
	{"exo", "exo.cue"},
	{"exohcl", "exo.hcl"},
	{"compose", "compose.yaml"},
	{"compose", "compose.yml"},
	{"compose", "docker-compose.yaml"},
	{"compose", "docker-compose.yml"},
	{"procfile", "Procfile"},
}

func (r *QueryResolver) findManifest(ctx context.Context, fs *FileSystemResolver, format *string) (*ManifestResolver, error) {
	for _, candidate := range manifestCandidates {
		if format != nil && *format != candidate.Format {
			continue
		}
		file, err := fs.file(ctx, "/"+candidate.Filename)
		if file == nil {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("searching for manifest: %w", err)
		}
		return &ManifestResolver{
			Format: candidate.Format,
			File:   file,
		}, nil
	}
	return nil, nil
}

func (r *ManifestResolver) Content() (string, error) {
	if r.File != nil {
		return r.File.Content()
	}
	return r.content, nil
}

func (r *ManifestResolver) Formatted() (string, error) {
	content, err := r.Content()
	if err != nil {
		return "", fmt.Errorf("resolving content: %w", err)
	}
	switch r.Format {
	case "exo":
		return cueutil.FormatString(content)
	default:
		// No-op for unsupported formats.
		return content, nil
	}
}

func (r *ManifestResolver) HostPath() *string {
	if r.File == nil {
		return nil
	}
	return &r.File.HostPath
}

func (r *MutationResolver) FormatManifest(ctx context.Context, args struct {
	Workspace string
	Format    *string
	Path      *string
}) (*VoidResolver, error) {
	workspace, err := r.workspaceByRef(ctx, &args.Workspace)
	if err := validateResolve("workspace", args.Workspace, workspace, err); err != nil {
		return nil, err
	}
	var manifest *ManifestResolver
	if args.Path == nil {
		if args.Format == nil {
			manifest, err = workspace.findManifest(ctx, args.Format)
		} else {
			manifest, err = workspace.Manifest(ctx)
		}
	} else {
		manifest, err = workspace.manifestByPath(ctx, *args.Path, args.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("resolving manifest: %w", err)
	}
	content, err := manifest.Formatted()
	if err != nil {
		return nil, err
	}
	hostPath := manifest.HostPath()
	if err := atomic.WriteFile(*hostPath, strings.NewReader(content)); err != nil {
		return nil, fmt.Errorf("writing manifest: %w", err)
	}
	return nil, nil
}
