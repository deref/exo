package file

import (
	"context"
	"errors"
	"os"
	"path"
)

func (c *Controller) Identify(ctx context.Context, m *Model) (string, error) {
	return path.Join("exo:/files", m.Path), nil
}

func (c *Controller) Create(ctx context.Context, m *Model) error {
	f, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(m.Content)
	return err
}

func (c *Controller) Read(ctx context.Context, m *Model) error {
	bs, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}
	m.Content = string(bs)
	return nil
}

func (c *Controller) Update(ctx context.Context, prev *Model, cur *Model) error {
	if prev.Path != cur.Path {
		return errors.New("moving file would change identity")
	}
	if prev.Content == cur.Content {
		return nil
	}
	return os.WriteFile(cur.Path, []byte(cur.Content), 0600)
}

func (c *Controller) Delete(ctx context.Context, m *Model) error {
	err := os.Remove(m.Path)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return err
}
