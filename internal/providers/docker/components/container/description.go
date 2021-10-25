package container

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
	dockerclient "github.com/docker/docker/client"
	"github.com/moby/moby/errdefs"
	"golang.org/x/sync/errgroup"
)

func GetProcessDescription(ctx context.Context, dockerClient *dockerclient.Client, component api.ComponentDescription) (api.ProcessDescription, error) {
	if component.Type != "container" {
		return api.ProcessDescription{}, fmt.Errorf("component not a container")
	}

	var state State
	if err := jsonutil.UnmarshalStringOrEmpty(component.State, &state); err != nil {
		return api.ProcessDescription{}, fmt.Errorf("unmarshalling container state: %v\n", err)
	}
	process := api.ProcessDescription{
		ID:       component.ID,
		Name:     component.Name,
		Provider: "docker",
	}

	containerInfo, err := dockerClient.ContainerInspect(ctx, state.ContainerID)
	if err != nil {
		// If there is an error inspecting the container, assume that this is
		// because the container hasn't been created yet and return the information
		// we already have.
		return process, nil
	}

	process.Running = containerInfo.State.Running

	startTime, err := time.Parse(time.RFC3339Nano, containerInfo.State.StartedAt)
	if err != nil {
		return api.ProcessDescription{}, fmt.Errorf("could not parse start time: %w", err)
	}
	createTime := startTime.UnixNano() / 1e6
	process.CreateTime = &createTime

	process.EnvVars = map[string]string{}
	for _, env := range containerInfo.Config.Env {
		decomposedEnv := strings.SplitN(env, "=", 2)
		process.EnvVars[decomposedEnv[0]] = decomposedEnv[1]
	}

	process.Ports = []uint32{}
	for p := range containerInfo.NetworkSettings.Ports {
		process.Ports = append(process.Ports, uint32(p.Int()))
	}

	var eg errgroup.Group
	eg.Go(func() error {
		statRequest, err := dockerClient.ContainerStats(ctx, state.ContainerID, false)
		if err != nil && !errdefs.IsConflict(err) { // ignore not running errors
			return fmt.Errorf("getting stats for container: %w", err)
		}
		defer statRequest.Body.Close()

		var containerStats docker.ContainerStats
		decoder := json.NewDecoder(statRequest.Body)
		if err := decoder.Decode(&containerStats); err != nil {
			return fmt.Errorf("could not unmarshal container stats: %s", err)
		}

		process.ResidentMemory = &containerStats.MemoryStats.Usage
		if containerStats.CPUStats.SystemCPUUsage != 0 {
			cpuPercent := float64(containerStats.CPUStats.CPUUsage.TotalUsage) / 1e9
			process.CPUPercent = &cpuPercent
		}
		return nil
	})

	eg.Go(func() error {
		topBody, err := dockerClient.ContainerTop(ctx, state.ContainerID, []string{})
		if err != nil && !errdefs.IsConflict(err) { // ignore not running errors
			return fmt.Errorf("running top in container: %w", err)
		}

		process.ChildrenExecutables = make([]string, len(topBody.Processes))
		for i, proc := range topBody.Processes {
			process.ChildrenExecutables[i] = proc[len(proc)-1]
		}
		return nil
	})

	// No context for this error since it's just collecting an already
	// contextualised error.
	err = eg.Wait()
	return process, err
}
