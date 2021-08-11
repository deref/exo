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
	if c.ImageID != "" {
		// TODO: When should we rebuild?
		return nil
	}
	if c.canBuild() {
		return c.buildImage(ctx)
	}
	inspection, _, err := c.Docker.ImageInspectWithRaw(ctx, c.Image)
	if docker.IsErrNotFound(err) {
		if err := c.pullImage(ctx); err != nil {
			return fmt.Errorf("pulling image: %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspecting image: %w", err)
	}
	c.ImageID = inspection.ID

	c.Command = make([]string, len(inspection.Config.Cmd))
	for i, cmd := range inspection.Config.Cmd {
		c.Command[i] = cmd
	}

	return nil
}

func (c *Container) pullImage(ctx context.Context) error {
	image, err := c.Docker.ImagePull(ctx, c.Image, types.ImagePullOptions{
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
