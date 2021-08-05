package osutil

import "syscall"

func SignalProcessGroup(pgid int, sig syscall.Signal) error {
	return syscall.Kill(-pgid, sig)
}

func KillProcessGroup(pgid int) error {
	return SignalProcessGroup(pgid, syscall.SIGKILL)
}
