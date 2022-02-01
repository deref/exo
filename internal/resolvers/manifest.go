package resolvers

import (
	"context"
	"fmt"
)

type ManifestResolver struct {
	Format string
	File   *FileResolver
}

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

func (r *QueryResolver) findManifest(ctx context.Context, fs *FileSystemResolver, format string) (*ManifestResolver, error) {
	for _, candidate := range manifestCandidates {
		if format != "" && format != candidate.Format {
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
			Format: format,
			File:   file,
		}, nil
	}
	return nil, nil
}

func (r *ManifestResolver) HostPath() string {
	return r.File.HostPath
}
