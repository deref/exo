package compose

import (
	"fmt"
	"strconv"

	"code.cloudfoundry.org/bytefmt"
)

type Memory int64

func (memory *Memory) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var memString string
	if err := unmarshal(&memString); err != nil {
		return err
	}
	memBytes, err := strconv.ParseInt(memString, 10, 64)
	if err == nil {
		*memory = Memory(memBytes)
		return nil
	}

	uMemBytes, err := bytefmt.ToBytes(memString)
	if err == nil {
		*memory = Memory(uMemBytes)
		return nil
	}

	return fmt.Errorf("unmarshaling memory value: %w", err)
}
