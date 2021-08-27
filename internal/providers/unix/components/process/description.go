package process

import (
	"context"
	"fmt"

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
	if err == nil {
		process.Running, err = proc.IsRunning()
		if err != nil {
			return api.ProcessDescription{}, err
		}

		memoryInfo, err := proc.MemoryInfo()
		if err != nil {
			return api.ProcessDescription{}, err
		}

		process.ResidentMemory = &memoryInfo.RSS

		connections, err := proc.Connections()
		if err != nil {
			return api.ProcessDescription{}, err
		}

		var ports []uint32
		for _, conn := range connections {
			if conn.Laddr.Port != 0 {
				ports = append(ports, conn.Laddr.Port)
			}
		}
		process.Ports = ports

		*process.CreateTime, err = proc.CreateTime()
		if err != nil {
			return api.ProcessDescription{}, err
		}

		children, err := proc.Children()
		if err == nil {
			var childrenExecutables []string
			for _, child := range children {
				exe, err := child.Exe()
				if err != nil {
					return api.ProcessDescription{}, err
				}
				childrenExecutables = append(childrenExecutables, exe)
			}
			process.ChildrenExecutables = childrenExecutables
		}

		*process.CPUPercent, err = proc.CPUPercent()
		if err != nil {
			return api.ProcessDescription{}, err
		}
	}

	return process, nil
}
