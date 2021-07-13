package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var child *os.Process
var varDir string

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fatalf("usage: %s <vardir> <command> <args...>", os.Args[0])
	}
	varDir = args[0]
	command := args[1]
	arguments := args[2:]

	cmd := exec.Command(command, arguments...)
	cmd.Env = os.Environ()

	// Connect pipes.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start child process.
	if err := cmd.Start(); err != nil {
		fatalf("starting %q: %v", command, err)
	}
	child = cmd.Process

	// Reporting child pid to stdout.
	if _, err := fmt.Println(child.Pid); err != nil {
		fatalf("reporting pid: %v", err)
	}

	// Forward signals to child.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for sig := range c {
			if err := cmd.Process.Signal(sig); err != nil {
				break
			}
		}
	}()

	// Proxy logs.
	go pipeToFifo("out", stdout)
	go pipeToFifo("err", stderr)

	// Wait for child process to exit.
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}
	if err != nil {
		fatalf("wait error: %v", err)
	}
}

func pipeToFifo(name string, r io.Reader) {
	b := bufio.NewReader(r)
	fifoPath := filepath.Join(varDir, name)
	if err := syscall.Mkfifo(fifoPath, 0600); err != nil {
		fatalf("making fifo %q: %v", fifoPath, err)
	}
	for {
		f, err := os.OpenFile(fifoPath, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			fatalf("opening fifo %q: %v", fifoPath, err)
		}
		for {
			line, isPrefix, err := b.ReadLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				fatalf("reading %s: %v", name, err)
			}
			// TODO: Do something better with lines that are too long.
			for isPrefix {
				// Skip remainder of line.
				line = append([]byte{}, line...)
				_, isPrefix, err = b.ReadLine()
				if err == io.EOF {
					return
				}
				if err != nil {
					fatalf("reading %s: %v", name, err)
				}
			}
			if _, err := f.Write(line); err != nil {
				fatalf("forwarding %s: %v", name, err)
			}
		}
	}
}

func fatalf(format string, v ...interface{}) {
	if child != nil {
		_ = child.Kill()
	}
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}
