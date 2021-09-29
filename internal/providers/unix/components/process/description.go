package process

import (
	"context"
	"fmt"

	psprocess "github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sync/errgroup"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/jsonutil"
)

func GetProcessDescription(ctx context.Context, component api.ComponentDescription) (api.ProcessDescription, error) {
	var state State
	if err := jsonutil.UnmarshalStringOrEmpty(component.State, &state); err != nil {
		return api.ProcessDescription{}, fmt.Errorf("unmarshalling container state: %v\n", err)
	}

	process := api.ProcessDescription{
		ID:       component.ID,
		Name:     component.Name,
		Provider: "unix",
		EnvVars:  state.FullEnvironment,
		Spec:     component.Spec,
	}

	proc, err := psprocess.NewProcess(int32(state.Pid))
	if err != nil {
		// Assume this has failed because the process isn't running.
		return process, nil
	}
	process.Running = true

	var eg errgroup.Group

	eg.Go(func() error {
		memoryInfo, err := proc.MemoryInfoWithContext(ctx)
		if err != nil {
			return fmt.Errorf("getting process memory information: %w", err)
		}
		process.ResidentMemory = &memoryInfo.RSS
		return nil
	})

	eg.Go(func() error {
		connections, err := proc.ConnectionsWithContext(ctx)
		if err != nil {
			return fmt.Errorf("getting process connections information: %w", err)
		}

		ports := []uint32{}
		for _, conn := range connections {
			if conn.Laddr.Port != 0 {
				ports = append(ports, conn.Laddr.Port)
			}
		}
		process.Ports = ports
		return nil
	})

	eg.Go(func() error {
		createTime, err := proc.CreateTimeWithContext(ctx)
		if err != nil {
			return fmt.Errorf("getting process createTime information: %w", err)
		}
		process.CreateTime = &createTime
		return nil
	})

	eg.Go(func() error {
		children, err := proc.ChildrenWithContext(ctx)
		if err != nil {
			// Assume that this has failed because the process doesn't have any
			// children.
			return nil
		}

		var childrenExecutables []string
		for _, child := range children {
			exe, err := child.Exe()
			if err != nil {
				continue
			}
			childrenExecutables = append(childrenExecutables, exe)
		}
		process.ChildrenExecutables = childrenExecutables
		return nil
	})

	eg.Go(func() error {
		cpuPercent, err := proc.CPUPercentWithContext(ctx)
		if err != nil {
			return fmt.Errorf("getting process cpuPercent information: %w", err)
		}
		process.CPUPercent = &cpuPercent
		return nil
	})

	err = eg.Wait()
	return process, err
}
