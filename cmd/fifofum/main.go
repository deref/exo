// Command fifofum execs a command with redirected stdio as given by arguments.
// This is used to enable redirecting to named pipes without having to pass
// file descriptors between processes.

package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	args := os.Args[1:]
	if len(args) < 4 {
		fatalf("usage: %s <stdin> <stdout> <stderr> <args...>", os.Args[0])
	}

	openr("stdin", args[0], syscall.Stdin)
	openw("stdout", args[1], syscall.Stdout)
	openw("stderr", args[2], syscall.Stderr)

	argv0 := args[3]
	argv := args[4:]
	envv := os.Environ()
	if err := syscall.Exec(argv0, argv, envv); err != nil {
		fatalf("executing %q: %w", argv0, err)
	}
}

func openr(role string, path string, replace int) {
	f, err := os.Open(path)
	if err != nil {
		fatalf("opening %q for %s: %w", path, role, err)
	}
	replaceFD(role, f, replace)
}

func openw(role string, path string, replace int) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		fatalf("opening %q for %s: %w", path, role, err)
	}
	replaceFD(role, f, replace)
}

func replaceFD(role string, from *os.File, to int) {
	fd := int(from.Fd())
	if err := syscall.Dup2(fd, to); err != nil {
		fatalf("replacing %s file descriptor %d with %d: %w", role, fd, to, err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
