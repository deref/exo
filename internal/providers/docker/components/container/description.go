package container

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	dockerclient "github.com/docker/docker/client"
)

func GetProcessDescription(ctx context.Context, dockerClient *dockerclient.Client, component api.ComponentDescription) (api.ProcessDescription, error) {
	logger := logging.CurrentLogger(ctx)
	if component.Type != "container" {
		return api.ProcessDescription{}, fmt.Errorf("component not a container")
	}

	var state State
	if err := jsonutil.UnmarshalString(component.State, &state); err != nil {
		return api.ProcessDescription{}, fmt.Errorf("unmarshalling container state: %v\n", err)
	}
	process := api.ProcessDescription{
		ID:       component.ID,
		Name:     component.Name,
		Provider: "docker",
	}

	containerInfo, err := dockerClient.ContainerInspect(ctx, state.ContainerID)
	if err != nil {
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

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		statRequest, err := dockerClient.ContainerStats(ctx, state.ContainerID, false)
		if err != nil {
			return
		}
		defer statRequest.Body.Close()

		var containerStats docker.ContainerStats
		decoder := json.NewDecoder(statRequest.Body)
		if err := decoder.Decode(&containerStats); err != nil {
			logger.Infof(fmt.Sprintf("could not unmarshal container stats: %s", err))
			return
		}

		process.ResidentMemory = &containerStats.MemoryStats.Usage
		if containerStats.CPUStats.SystemCPUUsage != 0 {
			cpuPercent := float64(containerStats.CPUStats.CPUUsage.TotalUsage) / 1e9
			process.CPUPercent = &cpuPercent
		}
	}()

	go func() {
		wg.Add(1)
		topBody, err := dockerClient.ContainerTop(ctx, state.ContainerID, []string{})
		if err != nil {
			return
		}

		process.ChildrenExecutables = make([]string, len(topBody.Processes))
		for i, proc := range topBody.Processes {
			process.ChildrenExecutables[i] = proc[len(proc)-1]
		}
	}()

	wg.Wait()

	return process, nil
}
