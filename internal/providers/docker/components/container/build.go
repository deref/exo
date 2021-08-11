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

	"github.com/deref/exo/internal/util/pathutil"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"golang.org/x/sync/errgroup"
)

func (c *Container) canBuild() bool {
	return c.Build.Context != ""
}

func (c *Container) buildImage(ctx context.Context) error {
	buildContext, buildContextWriter := io.Pipe()
	defer buildContextWriter.Close()
	var eg errgroup.Group

	eg.Go(func() error {
		contextPath := filepath.Join(c.WorkspaceRoot, c.Build.Context)
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
			Isolation: container.Isolation(c.Build.Isolation),
			//CPUSetCPUs     string
			//CPUSetMems     string
			//CPUShares      int64
			//CPUQuota       int64
			//CPUPeriod      int64
			//Memory         int64
			//MemorySwap     int64
			//CgroupParent   string
			//NetworkMode    string
			ShmSize:    int64(c.Build.ShmSize),
			Dockerfile: c.Build.Dockerfile,
			//Ulimits        []*units.Ulimit
			//// BuildArgs needs to be a *string instead of just a string so that
			//// we can tell the difference between "" (empty string) and no value
			//// at all (nil). See the parsing of buildArgs in
			//// api/server/router/build/build_routes.go for even more info.
			BuildArgs: c.Build.Args,
			//AuthConfigs map[string]AuthConfig
			//Context     io.Reader
			Labels: c.Build.Labels.WithoutNils(),
			//// squash the resulting image's layers to the parent
			//// preserves the original image and creates a new one from the parent with all
			//// the changes applied to a single layer
			//Squash bool
			// CacheFrom specifies images that are used for matching cache. Images
			// specified here do not need to have a valid parent chain to match cache.
			CacheFrom: c.Build.CacheFrom,
			//SecurityOpt []string
			ExtraHosts: c.Build.ExtraHosts,
			Target:     c.Build.Target,
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
		fmt.Println("before build")
		resp, err := c.Docker.ImageBuild(ctx, buildContext, opts)
		fmt.Println("after build")
		if resp.Body != nil {
			defer resp.Body.Close()
			// TODO [DOCKER_PROGRESS]: Capture progress.
			//_, _ = io.Copy(os.Stdout, resp.Body)
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				var d struct {
					ErrorDetail struct {
						Code    int    `json:"code"`
						Message string `json:"message"`
					} `json:"errorDetail"`
					Aux struct {
						ID string `json:"ID"`
					} `json:"aux"`
				}
				if err := json.Unmarshal(scanner.Bytes(), &d); err != nil {
					return fmt.Errorf("failed to unmarshal docker build log: %w", err)
				}
				if d.ErrorDetail.Message != "" {
					return fmt.Errorf("docker build error: " + d.ErrorDetail.Message)
				}
				if strings.HasPrefix(d.Aux.ID, "sha256:") {
					c.State.ImageID = d.Aux.ID
				}
			}
		}
		if err != nil {
			return err
		}
		if c.State.ImageID == "" {
			return fmt.Errorf("did not build an image")
		}
		return nil
	})

	return eg.Wait()
}

func tarBuildContext(w io.WriteCloser, root string) (err error) {
	defer func() {
		err2 := w.Close()
		if err == nil {
			err = err2
		}
	}()

	tw := tar.NewWriter(w)

	filepath.Walk(root, func(file string, info os.FileInfo, err error) error {
		fmt.Println("taring", file)
		// Generate and write file header.
		header, err := tar.FileInfoHeader(info, file)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(root, filepath.ToSlash(file))
		if err != nil {
			return err
		}
		fmt.Println(">>", header.Name)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If not a directory, write file content.
		if !info.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	fmt.Println("tarred")
	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}
