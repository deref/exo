package osutil

import (
	"os"
	"syscall"
)

func IsValidPid(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
