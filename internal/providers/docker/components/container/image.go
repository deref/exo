package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func (c *Container) ensureImage(ctx context.Context) error {
	if c.State.Image.ID != "" {
		// TODO: When should we rebuild?
		return nil
	}
	if c.canBuild() {
		return c.buildImage(ctx)
	}

	inspection, _, err := c.Docker.ImageInspectWithRaw(ctx, c.Spec.Image)
	if docker.IsErrNotFound(err) {
		if err := c.pullImage(ctx); err != nil {
			return fmt.Errorf("pulling image: %w", err)
		}
		inspection, _, err = c.Docker.ImageInspectWithRaw(ctx, c.Spec.Image)
	}
	if err != nil {
		return fmt.Errorf("inspecting image: %w", err)
	}

	c.State.Image.ID = inspection.ID
	c.State.Image.Command = inspection.Config.Cmd
	c.State.Image.WorkingDir = inspection.Config.WorkingDir
	c.State.Image.Entrypoint = inspection.Config.Entrypoint
	c.State.Image.Shell = inspection.Config.Shell
	if len(c.State.Image.Shell) == 0 {
		if inspection.Os == "linux" {
			c.State.Image.Shell = []string{"/bin/sh", "-c"}
		} else {
			// For Windows — this is untested but it is what docker does.
			c.State.Image.Shell = []string{"cmd", "/S", "/C"}
		}

	}
	return nil
}

func (c *Container) pullImage(ctx context.Context) error {
	image, err := c.Docker.ImagePull(ctx, c.Spec.Image, types.ImagePullOptions{
		//All           bool
		//RegistryAuth  string // RegistryAuth is the base64 encoded credentials for the registry
		//PrivilegeFunc RequestPrivilegeFunc
		//Platform      string
	})
	if image != nil {
		// TODO [DOCKER_PROGRESS]: Report progress somehow.
		_, _ = io.Copy(os.Stdout, image)
		_ = image.Close()
	}
	return err
}
