// SEE NOTE: [IMAGE_SUBCOMPONENT].

package container

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deref/exo/internal/providers/docker/components/image"
	"github.com/deref/exo/internal/task"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func (c *Container) ensureImage(ctx context.Context, spec *Spec) error {
	if c.State.Image.ID != "" && spec.PullPolicy.Value != "build" {
		return nil
	}

	imageSpec := &image.Spec{
		Platform: spec.Platform.Value,
		Build:    spec.Build,
	}

	var inspection types.ImageInspect
	var err error
	if c.canBuild(imageSpec) {
		if err := c.buildImage(ctx, imageSpec); err != nil {
			return fmt.Errorf("building image: %w", err)
		}
		inspection, _, err = c.Docker.ImageInspectWithRaw(ctx, c.State.Image.ID)
		if err != nil {
			return fmt.Errorf("inspecting built image: %w", err)
		}
	} else {
		if spec.PullPolicy.Value != "always" {
			inspection, _, err = c.Docker.ImageInspectWithRaw(ctx, spec.Image.Value)
			if docker.IsErrNotFound(err) {
				if spec.PullPolicy.Value == "never" {
					return fmt.Errorf("pull policy for %q set to \"never\", no image %q found in local cache, and no build specification provided", c.ComponentName, spec.Image)
				}
			} else if err != nil {
				return fmt.Errorf("inspecting image: %w", err)
			}
		}
		if inspection.ID == "" {
			if err := c.pullImage(ctx, spec); err != nil {
				return fmt.Errorf("pulling image: %w", err)
			}
			inspection, _, err = c.Docker.ImageInspectWithRaw(ctx, spec.Image.Value)
			if err != nil {
				return fmt.Errorf("inspecting pulled image: %w", err)
			}
		}
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

type dockerPullStatus struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Progress       string `json:"string"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func (c *Container) pullImage(ctx context.Context, spec *Spec) error {
	pullTask := task.CurrentTask(ctx)
	if pullTask == nil {
		panic("No build task")
	}

	image, err := c.Docker.ImagePull(ctx, spec.Image.Value, types.ImagePullOptions{
		//All           bool
		//RegistryAuth  string // RegistryAuth is the base64 encoded credentials for the registry
		//PrivilegeFunc RequestPrivilegeFunc
		//Platform      string
	})
	if image != nil {
		defer image.Close()
		decoder := json.NewDecoder(image)
		for decoder.More() {
			var status dockerPullStatus
			if err := decoder.Decode(&status); err != nil {
				pullTask.Fail(fmt.Errorf("decoding image pull status: %w", err))
				break
			}
			pullTask.ReportProgress(status.ProgressDetail.Current, status.ProgressDetail.Total)
			pullTask.ReportMessage(status.Status)
		}
	}
	return err
}
