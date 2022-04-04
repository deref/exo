// This file provides a controller for typical operating system files.  This
// component type is not super useful in practice, but it's very useful for
// testing the controller system.

package os

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/sdk"
)

func NewFileController(svc api.Service) *sdk.ResourceComponentController {
	return sdk.NewResourceComponentController[FileModel](svc, &FileController{})
}

type FileController struct{}

type FileModel struct {
	FileSpec
	FileState
}

type FileSpec struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type FileState struct {
	Size     *string          `json:"size,omitempty"` // String for int64 support.
	Modified *scalars.Instant `json:"modified,omitempty"`
}

func (ctrl *FileController) IdentifyResource(ctx context.Context, cfg *sdk.ResourceConfig, m *FileModel) (string, error) {
	return path.Join("exo:/files", m.Path), nil
}

func (ctrl *FileController) CreateResource(ctx context.Context, cfg *sdk.ResourceConfig, m *FileModel) error {
	return ctrl.write(m, os.O_CREATE)
}

func (ctrl *FileController) ReadResource(ctx context.Context, cfg *sdk.ResourceConfig, m *FileModel) error {
	f, err := os.Open(m.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ctrl.stat(f, m); err != nil {
		return err
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	m.Content = string(bs)
	return nil
}

func (ctrl *FileController) UpdateResource(ctx context.Context, cfg *sdk.ResourceConfig, prev *FileModel, next *FileModel) error {
	if prev.Path != next.Path {
		return sdk.NewNotImplementedErrorf("moving file would change identity")
	}
	if prev.Content == next.Content {
		return nil
	}
	return ctrl.write(next, os.O_TRUNC)
}

func (ctrl *FileController) write(m *FileModel, flag int) error {
	f, err := os.OpenFile(m.Path, flag|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(m.Content)
	if err != nil {
		return err
	}

	return ctrl.stat(f, m)
}

func (ctrl *FileController) stat(f *os.File, m *FileModel) error {
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := strconv.FormatInt(stat.Size(), 10)
	m.Size = &size

	modified := scalars.GoTimeToInstant(stat.ModTime())
	m.Modified = &modified

	return nil
}

func (ctrl *FileController) ShutdownResource(ctx context.Context, cfg *sdk.ResourceConfig, m *FileModel) error {
	return nil
}

func (ctrl *FileController) DeleteResource(ctx context.Context, cfg *sdk.ResourceConfig, m *FileModel) error {
	err := os.Remove(m.Path)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return err
}
