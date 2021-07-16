package osutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
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

func ReadPid(path string) int {
	bs, _ := ioutil.ReadFile(path)
	pid, _ := strconv.Atoi(string(bytes.TrimSpace(bs)))
	return pid
}

func WritePid(path string, pid int) error {
	return ioutil.WriteFile(path, []byte(strconv.Itoa(pid)+"\n"), 0600)
}
