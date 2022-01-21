package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
)

type Controller struct{}

func (c *Controller) Identify(ctx context.Context, m *Model) (string, error) {
	// TODO: Validate host id & path.
	return path.Join(fmt.Sprintf("exo:hosts/%s/files", m.Path)), nil
}

func (c *Controller) Create(ctx context.Context, m *Model) error {
	// TODO: Verify current host id.
	f, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(m.Content)
	return err
}

func (c *Controller) Read(ctx context.Context, m *Model) error {
	// TODO: Verify current host id.
	bs, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}
	m.Content = string(bs)
	return nil
}

func (c *Controller) Update(ctx context.Context, prev *Model, cur *Model) error {
	if prev.HostID != cur.HostID || prev.Path != cur.Path {
		return errors.New("moving file would change identity")
	}
	// TODO: Verify current host id.
	if prev.Content == cur.Content {
		return nil
	}
	return os.WriteFile(cur.Path, []byte(cur.Content), 0600)
}
