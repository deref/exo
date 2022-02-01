package resolvers

import (
	"context"
	"errors"
	"path"

	"github.com/deref/exo/internal/util/pathutil"
)

type FileSystemResolver struct {
	hostPath string
}

func (r *FileSystemResolver) Root(ctx context.Context) (*FileResolver, error) {
	file, err := r.file(ctx, "/")
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("target directory does not exist")
	}
	return file, nil
}

func (r *FileSystemResolver) File(ctx context.Context, args struct {
	Path string
}) (*FileResolver, error) {
	return r.file(ctx, args.Path)
}

func (r *FileSystemResolver) file(ctx context.Context, exposedPath string) (*FileResolver, error) {
	if exposedPath == "" || exposedPath[0] != '/' {
		return nil, errors.New("invalid absolute path")
	}
	hostPath := path.Join(r.hostPath, exposedPath)
	if !pathutil.HasPathPrefix(hostPath, r.hostPath) {
		return nil, errors.New("path escapes filesystem root")
	}
	return resolveFile(exposedPath, hostPath)
}
