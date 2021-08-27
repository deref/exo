package process

import (
	"context"
	"fmt"
	"sync"

	psprocess "github.com/shirou/gopsutil/v3/process"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/jsonutil"
)

func GetProcessDescription(ctx context.Context, component api.ComponentDescription) (api.ProcessDescription, error) {
	var state State
	if err := jsonutil.UnmarshalString(component.State, &state); err != nil {
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
		return api.ProcessDescription{}, fmt.Errorf("getting process: %w", err)
	}

	var wg sync.WaitGroup
	work := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	work(func() {
		running, err := proc.IsRunningWithContext(ctx)
		process.Running = err == nil && running
	})

	work(func() {
		if memoryInfo, err := proc.MemoryInfoWithContext(ctx); err == nil {
			process.ResidentMemory = &memoryInfo.RSS
		}
	})

	work(func() {
		connections, err := proc.ConnectionsWithContext(ctx)
		if err != nil {
			return
		}

		var ports []uint32
		for _, conn := range connections {
			if conn.Laddr.Port != 0 {
				ports = append(ports, conn.Laddr.Port)
			}
		}
		process.Ports = ports

	})

	work(func() {
		if createTime, err := proc.CreateTimeWithContext(ctx); err == nil {
			process.CreateTime = &createTime
		}
	})

	work(func() {
		children, err := proc.ChildrenWithContext(ctx)
		if err != nil {
			return
		}

		var childrenExecutables []string
		for _, child := range children {
			exe, err := child.Exe()
			if err != nil {
				return
			}
			childrenExecutables = append(childrenExecutables, exe)
		}
		process.ChildrenExecutables = childrenExecutables

	})

	work(func() {
		if cpuPercent, err := proc.CPUPercentWithContext(ctx); err == nil {
			process.CPUPercent = &cpuPercent
		}
	})

	wg.Wait()

	return process, nil
}
