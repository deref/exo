package resolvers

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

type FileResolver struct {
	Path     string
	hostPath string
	info     fs.FileInfo
}

func openFileResolver(exposedPath, hostPath string) (*FileResolver, error) {
	r := FileResolver{
		Path:     exposedPath,
		hostPath: hostPath,
	}
	f, err := os.Open(hostPath)
	if f != nil {
		defer f.Close()
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err == nil {
		r.info, err = f.Stat()
		if err != nil {
			return nil, err
		}
	}
	return &r, nil
}

func (r *FileResolver) Name() string {
	return r.info.Name()
}

func (r *FileResolver) IsDirectory() bool {
	return r.info.IsDir()
}

func (r *FileResolver) Size() float64 {
	return float64(r.info.Size())
}

func (r *FileResolver) Content() (string, error) {
	bs, err := ioutil.ReadFile(r.hostPath)
	return string(bs), err
}

func (r *FileResolver) Children() ([]*FileResolver, error) {
	if !r.info.IsDir() {
		return nil, nil
	}
	f, err := os.Open(r.hostPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	entries, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	res := make([]*FileResolver, len(entries))
	for i, entry := range entries {
		name := entry.Name()
		exposedPath := path.Join(r.Path, name)
		hostPath := path.Join(r.hostPath, name)
		res[i], err = openFileResolver(exposedPath, hostPath)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
