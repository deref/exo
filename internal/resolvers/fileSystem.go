package resolvers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/deref/exo/internal/util/pathutil"
)

type FileSystemResolver struct {
	hostPath string
}

func (r *RootResolver) FileSystem() *FileSystemResolver {
	return r.fileSystemByHostPath("/") // TODO: Windows.
}

func (r *RootResolver) fileSystemByHostPath(hostPath string) *FileSystemResolver {
	return &FileSystemResolver{
		hostPath: hostPath,
	}
}

func (r *FileSystemResolver) HomePath() (string, error) {
	if r.hostPath == "/" {
		homePath, err := os.UserHomeDir()
		if err == nil {
			homePath = withDirectorySuffix(homePath)
		}
		return homePath, err
	} else {
		return "/", nil
	}
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

func (r *FileSystemResolver) FileOrHome(ctx context.Context, args struct {
	Path *string
}) (*FileResolver, error) {
	if args.Path != nil {
		return r.file(ctx, *args.Path)
	}
	homePath, err := r.HomePath()
	if err != nil {
		return nil, fmt.Errorf("resolving home: %w", err)
	}
	return r.file(ctx, homePath)
}

func (r *FileSystemResolver) file(ctx context.Context, exposedPath string) (*FileResolver, error) {
	if exposedPath == "" || exposedPath[0] != '/' {
		return nil, errors.New("invalid absolute path")
	}
	hostPath := path.Join(r.hostPath, exposedPath)
	if !pathutil.HasPathPrefix(hostPath, r.hostPath) {
		return nil, errors.New("path escapes filesystem root")
	}
	file := &FileResolver{
		Path:     exposedPath,
		HostPath: hostPath,
	}
	err := file.init()
	if os.IsNotExist(file.err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}
