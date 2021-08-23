package container

import (
	"archive/tar"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/util/pathutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/moby/moby/builder/dockerignore"
	"github.com/moby/moby/pkg/fileutils"
	"golang.org/x/sync/errgroup"
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
	buildContext, buildContextWriter := io.Pipe()
	defer buildContextWriter.Close()
	var eg errgroup.Group

	buildTask := task.CurrentTask(ctx)
	if buildTask == nil {
		panic("No build task")
	}

	spec := c.Spec

	eg.Go(func() error {
		contextPath := filepath.Join(c.WorkspaceRoot, spec.Build.Context)
		if !pathutil.HasFilePathPrefix(contextPath, c.WorkspaceRoot) {
			return errors.New("docker container build context path must be in exo workspace root")
		}
		if err := tarBuildContext(buildContextWriter, contextPath); err != nil {
			return fmt.Errorf("tarring build context: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
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
	})

	return eg.Wait()
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

func tarBuildContext(w io.WriteCloser, root string) (err error) {
	defer func() {
		err2 := w.Close()
		if err == nil {
			err = err2
		}
	}()

	ignore, err := createIgnoreMatcher(root)
	if err != nil {
		return fmt.Errorf("reading dockerignore file: %w", err)
	}

	tw := tar.NewWriter(w)

	filepath.Walk(root, func(file string, info os.FileInfo, err error) error {
		// Skip ignored files.
		name := info.Name()
		ignored, err := ignore.Matches(name)
		if err != nil {
			// This is expected to be unreachable.
			// SEE NOTE: [LAZY_PATTERN_MATCHER].
			return fmt.Errorf("matching %q: %w", name, err)
		}
		if ignored {
			return nil
		}

		// Generate and write file header.
		header, err := tar.FileInfoHeader(info, file)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(root, filepath.ToSlash(file))
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If not a directory, write file content.
		if !info.IsDir() {
			data, err := os.OpenFile(file, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}

func createIgnoreMatcher(buildContext string) (*fileutils.PatternMatcher, error) {
	patterns := []string{}

	// Read additional patterns from ignore file.
	ignorePath := filepath.Join(buildContext, ".dockerignore")
	file, err := os.Open(ignorePath)
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if file != nil {
		defer file.Close()

		morePatterns, err := dockerignore.ReadAll(file)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, morePatterns...)
	}

	// Build a pattern matcher.
	matcher, err := fileutils.NewPatternMatcher(patterns)
	if err != nil {
		return nil, err
	}
	// NOTE [LAZY_PATTERN_MATCHER]: PatternMatcher compiles patterns lazily. If
	// compilation of an individual pattern fails, the returned error does not
	// not report which. Since the Matches method always visits every pattern, we
	// can invoke it here to force pattern compilation eagerly. This won't help
	// identifying which pattern failed to compile, but it will at least prevent
	// the error from being reported confusingly in the context of a particular
	// match operation.
	if _, err := matcher.Matches(""); err != nil {
		return nil, err
	}
	return matcher, nil
}
