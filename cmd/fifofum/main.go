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

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/logcol/api"
)

var child *os.Process
var varDir string

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fatalf(`usage: %s <vardir> <command> <args...>

fifofum executes and supervises the given command. If successful, the child
pid is written to stdout and two fifo files are created in vardir: out and err.
The corresponding stdio streams will be proxied from the supervised process to
those fifos.`, os.Args[0])
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

	// Handle signals.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGCHLD)
	go func() {
		for sig := range c {
			switch sig {
			// Forward signals to child.
			case os.Interrupt, os.Kill:
				if err := cmd.Process.Signal(sig); err != nil {
					break
				}
			// Exit when child exits.
			case syscall.SIGCHLD:
				os.Exit(1)
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
	b := bufio.NewReaderSize(r, api.MaxMessageSize)
	fifoPath := filepath.Join(varDir, name) + ".fifo"
	if err := syscall.Mkfifo(fifoPath, 0600); err != nil && !os.IsExist(err) {
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
			line = append(line, '\n')
			_, _ = f.Write(line)
		}
	}
}

func fatalf(format string, v ...interface{}) {
	if child != nil {
		_ = child.Kill()
	}
	cmdutil.Fatalf(format, v...)
}
