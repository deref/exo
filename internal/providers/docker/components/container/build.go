package container

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/util/pathutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/moby/moby/builder/dockerignore"
	"github.com/moby/moby/pkg/archive"
	"github.com/moby/moby/pkg/fileutils"
)

func (c *Container) Build(ctx context.Context, input *api.BuildInput) (*api.BuildOutput, error) {
	if c.canBuild() {
		if err := c.buildImage(ctx); err != nil {
			return nil, err
		}
	}
	return &api.BuildOutput{}, nil
}

func (c *Container) canBuild() bool {
	return c.Spec.Build.Context != ""
}

func (c *Container) buildImage(ctx context.Context) error {
	buildTask := task.CurrentTask(ctx)
	if buildTask == nil {
		panic("No build task")
	}

	spec := c.Spec

	contextPath := filepath.Join(c.WorkspaceRoot, spec.Build.Context)
	if !pathutil.HasFilePathPrefix(contextPath, c.WorkspaceRoot) {
		return errors.New("docker container build context path must be in exo workspace root")
	}
	dockerfile := c.Spec.Build.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	buildContext, err := getArchive(contextPath, dockerfile)
	if err != nil {
		return fmt.Errorf("getting build context: %w", err)
	}

	opts := types.ImageBuildOptions{
		//Tags           []string
		//SuppressOutput bool
		//RemoteContext  string
		//NoCache        bool
		//Remove         bool
		//ForceRemove    bool
		//PullParent     bool
		Isolation: container.Isolation(spec.Build.Isolation),
		//CPUSetCPUs     string
		//CPUSetMems     string
		//CPUShares      int64
		//CPUQuota       int64
		//CPUPeriod      int64
		//Memory         int64
		//MemorySwap     int64
		//CgroupParent   string
		//NetworkMode    string
		ShmSize:    int64(spec.Build.ShmSize),
		Dockerfile: spec.Build.Dockerfile,
		//Ulimits        []*units.Ulimit
		//// BuildArgs needs to be a *string instead of just a string so that
		//// we can tell the difference between "" (empty string) and no value
		//// at all (nil). See the parsing of buildArgs in
		//// api/server/router/build/build_routes.go for even more info.
		BuildArgs: spec.Build.Args,
		//AuthConfigs map[string]AuthConfig
		//Context     io.Reader
		Labels: spec.Build.Labels.WithoutNils(),
		//// squash the resulting image's layers to the parent
		//// preserves the original image and creates a new one from the parent with all
		//// the changes applied to a single layer
		//Squash bool
		// CacheFrom specifies images that are used for matching cache. Images
		// specified here do not need to have a valid parent chain to match cache.
		CacheFrom: spec.Build.CacheFrom,
		//SecurityOpt []string
		ExtraHosts: spec.Build.ExtraHosts,
		Target:     spec.Build.Target,
		//SessionID   string
		//Platform    string
		//// Version specifies the version of the unerlying builder to use
		//Version BuilderVersion
		//// BuildID is an optional identifier that can be passed together with the
		//// build request. The same identifier can be used to gracefully cancel the
		//// build with the cancel request.
		//BuildID string
		//// Outputs defines configurations for exporting build results. Only supported
		//// in BuildKit mode
		//Outputs []ImageBuildOutput
	}
	resp, err := c.Docker.ImageBuild(ctx, buildContext, opts)
	if resp.Body != nil {
		defer resp.Body.Close()

		subtasks := make(map[string]*task.Task)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var event buildEvent
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				return fmt.Errorf("failed to unmarshal docker build log: %w", err)
			}

			if event.ErrorDetail.Message != "" {
				// TODO: Report error code too.
				return fmt.Errorf("docker build error: " + event.ErrorDetail.Message)
			}

			if event.ID != "" {
				var subtask *task.Task
				if subtask == nil && event.Status == "Pulling fs layer" {
					subtask = buildTask.StartChild("layer " + event.ID)
					subtasks[event.ID] = subtask
				} else {
					subtask = subtasks[event.ID]
				}
				if subtask != nil {
					if event.Status != "" {
						subtask.ReportMessage(event.Status)
					}
					if event.Status == "Pull complete" {
						_ = subtask.Finish()
					}
					if event.ProgressDetail.Total > 0 {
						subtask.ReportProgress(event.ProgressDetail.Current, event.ProgressDetail.Total)
					}
				}
			}

			if event.Stream != "" {
				message := strings.TrimSpace(event.Stream)
				buildTask.ReportMessage(message)
			}

			if strings.HasPrefix(event.Aux.ID, "sha256:") {
				c.State.Image.ID = event.Aux.ID
			}
		}
	}
	if err != nil {
		return err
	}
	if c.State.Image.ID == "" {
		return fmt.Errorf("did not build an image")
	}
	return nil

}

// See <github.com/docker/docker/pkg/jsonmessage>.
type buildEvent struct {
	// Image tag, layer IDs, and other identifiers, contingent on content of "status".
	ID string `json:"id"`
	// Log message, with trailing newline.
	Stream string `json:"stream"`
	// Ad-hoc event type with some payload data too.
	Status string `json:"status"`
	// Present for some statuses, such as "Downloading", and "Extracting".
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
	// A rendered progress bar. We prefer "progressDetail".
	Progress string `json:"progress"`
	// Non-zero code and/or non-empty message when something has gone wrong.
	ErrorDetail struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errorDetail"`
	// Reports a built ID.
	Aux struct {
		ID string `json:"ID"`
	} `json:"aux"`
}

func getArchive(contextDir, relDockerfile string) (io.ReadCloser, error) {
	var err error

	// Convert dockerfile name to a platform-independent one.
	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	f, err := os.Open(filepath.Join(contextDir, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var excludes []string
	if err == nil {
		excludes, err = dockerignore.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}

	// If .dockerignore mentions .dockerignore or the Dockerfile
	// then make sure we send both files over to the daemon
	// because Dockerfile is, obviously, needed no matter what, and
	// .dockerignore is needed to know if either one needs to be
	// removed. The daemon will remove them for us, if needed, after it
	// parses the Dockerfile.
	var includes = []string{"."}
	keepThem1, err := fileutils.Matches(".dockerignore", excludes)
	if err != nil {
		return nil, fmt.Errorf("matching .dockerignore: %w", err)
	}
	keepThem2, err := fileutils.Matches(relDockerfile, excludes)
	if err != nil {
		return nil, fmt.Errorf("matching Dockerfile: %w", err)
	}
	if keepThem1 || keepThem2 {
		includes = append(includes, ".dockerignore", relDockerfile)
	}

	return archive.TarWithOptions(contextDir, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: excludes,
		IncludeFiles:    includes,
	})
}
