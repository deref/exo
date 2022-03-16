package resolvers

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FileResolver struct {
	Path        string
	HostPath    string
	IsDirectory bool

	info fs.FileInfo
	err  error
}

// Initializes resolver by stat-ing the file.  IsDirectory may be preset with
// knowledge from readdir on the parent, or may be reset by a successful stat.
// There are various consistency issues and race conditions around readdir,
// open, and stat during resolving. A file resolver is still valid even if it's
// err field is set, but some usages may choose to discard the resolver. For
// example, a file-not-exists error may mean replacing the resolver with nil.
// Additionally, directories will have their paths updated to always include a
// trailing separator.
func (r *FileResolver) init() error {
	var f *os.File
	f, r.err = os.Open(r.HostPath)
	if r.err == nil {
		defer f.Close()
		r.info, r.err = f.Stat()
		if r.info != nil {
			r.IsDirectory = r.info.IsDir()
		}
	}
	if r.IsDirectory {
		r.Path = withDirectorySuffix(r.Path)
		r.HostPath = withDirectorySuffix(r.HostPath)
	}
	return r.err
}

func withDirectorySuffix(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

func (r *FileResolver) Name() string {
	return path.Base(r.Path)
}

func (r *FileResolver) Size() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return float64(r.info.Size()), nil
}

func (r *FileResolver) Content() (string, error) {
	bs, err := ioutil.ReadFile(r.HostPath)
	return string(bs), err
}

func (r *FileResolver) ParentPath() *string {
	dir := path.Dir(r.Path)
	if dir == "/" {
		return nil
	}
	return &dir
}

func (r *FileResolver) Children() ([]*FileResolver, error) {
	if !r.info.IsDir() {
		return nil, nil
	}
	f, err := os.Open(r.HostPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	entries, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	children := make([]*FileResolver, len(entries))
	for i, entry := range entries {
		name := entry.Name()
		exposedPath := path.Join(r.Path, name)
		hostPath := path.Join(r.HostPath, name)
		child := &FileResolver{
			Path:        exposedPath,
			HostPath:    hostPath,
			IsDirectory: entry.IsDir(),
		}
		child.init()
		children[i] = child
	}
	return children, nil
}
