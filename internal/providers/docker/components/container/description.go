package container

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
	dockerclient "github.com/docker/docker/client"
)

func GetProcessDescription(ctx context.Context, dockerClient *dockerclient.Client, component api.ComponentDescription) (api.ProcessDescription, error) {
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
	if process.Running {
		statRequest, err := dockerClient.ContainerStats(ctx, state.ContainerID, false)
		if err != nil {
			return api.ProcessDescription{}, fmt.Errorf("could not get container stats: %w", err)
		}
		defer statRequest.Body.Close()

		scanner := bufio.NewScanner(statRequest.Body)
		didScan := scanner.Scan()
		if err = scanner.Err(); err != nil {
			return api.ProcessDescription{}, fmt.Errorf("error scanning container stats: %w", err)
		}

		if !didScan {
			return api.ProcessDescription{}, fmt.Errorf("could not read container stats")
		}

		var containerStats docker.ContainerStats
		err = json.Unmarshal(scanner.Bytes(), &containerStats)
		if err != nil {
			return api.ProcessDescription{}, fmt.Errorf("could not unmarshal container stats: %w", err)
		}

		process.ResidentMemory = containerStats.MemoryStats.Usage
		layout := "2006-01-02T15:04:05.000000000Z"
		startTime, err := time.Parse(layout, containerInfo.State.StartedAt)
		if err != nil {
			return api.ProcessDescription{}, fmt.Errorf("could not parse start time: %w", err)
		}

		process.CreateTime = startTime.UnixNano() / 1e6

		if containerStats.CPUStats.SystemCPUUsage != 0 {
			process.CPUPercent = float64(containerStats.CPUStats.CPUUsage.TotalUsage) / 1e9
		}

		process.EnvVars = map[string]string{}
		for _, env := range containerInfo.Config.Env {
			decomposedEnv := strings.SplitN(env, "=", 2)
			process.EnvVars[decomposedEnv[0]] = decomposedEnv[1]
		}

		process.Ports = []uint32{}
		for p := range containerInfo.NetworkSettings.Ports {
			process.Ports = append(process.Ports, uint32(p.Int()))
		}

		topBody, err := dockerClient.ContainerTop(ctx, state.ContainerID, []string{})
		if err != nil {
			return api.ProcessDescription{}, fmt.Errorf("could not top container: %w", err)
		}

		process.ChildrenExecutables = []string{}
		for _, proc := range topBody.Processes {
			process.ChildrenExecutables = append(process.ChildrenExecutables, proc[len(proc)-1])
		}

	}

	return process, nil
}
