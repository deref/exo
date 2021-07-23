package sysutil

import (
	"syscall"
)

func Dup2(from int, to int) error {
	flags := 0
	return syscall.Dup3(from, to, flags)
}
