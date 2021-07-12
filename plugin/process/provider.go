package process

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/deref/exo/api"
	"github.com/deref/exo/gensym"
	"github.com/mitchellh/mapstructure"
)

type Provider struct {
	ProjectDir string
	VarDir     string
}

type spec struct {
	Directory string
	Command   string
	Arguments []string
}

func (provider *Provider) Create(ctx context.Context, input *api.CreateInput) (*api.CreateOutput, error) {
	var spec spec
	if err := mapstructure.Decode(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("decoding mapstructure: %w", err)
	}

	var stdin, stdout, stderr *os.File
	var proc *os.Process
	abort := func() {
		if stdin != nil {
			_ = stdin.Close()
		}
		if stdout != nil {
			// TODO: Delete file.
			_ = stdout.Close()
		}
		if stderr != nil {
			// TODO: Delete file.
			_ = stderr.Close()
		}
		if proc != nil {
			proc.Kill()
		}
	}

	err := os.Mkdir(provider.VarDir, 0700)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("creating var directory: %w", err)
	}

	globalID := gensym.Base32()

	procDir := filepath.Join(provider.VarDir, globalID)
	if err := os.Mkdir(procDir, 0700); err != nil {
		return nil, fmt.Errorf("creating proc directory: %w", err)
	}

	stdin, err = os.Open("/dev/null")
	if err != nil {
		abort()
		return nil, fmt.Errorf("opening /dev/null: %w", err)
	}

	if err = mkfifo(procDir, "out"); err != nil {
		abort()
		return nil, err
	}
	if err := mkfifo(procDir, "err"); err != nil {
		abort()
		return nil, err
	}

	directory := spec.Directory
	if directory == "" {
		directory = provider.ProjectDir
	}

	env := []string{} // TODO: Fill from config.

	// XXX use fifofum
	proc, err = os.StartProcess(spec.Command, spec.Arguments, &os.ProcAttr{
		Dir:   directory,
		Files: []*os.File{stdin, stdout, stderr},
		Env:   addCriticalEnv(env),
		Sys:   nil,
	})
	if err != nil {
		abort()
		return nil, fmt.Errorf("starting %q: %w", spec.Command, err)
	}

	// Write pid file.
	if err := ioutil.WriteFile(filepath.Join(procDir, "pid"), []byte(globalID+"\n"), 0600); err != nil {
		abort()
		return nil, fmt.Errorf("writing pid: %w", err)
	}

	var output api.CreateOutput
	output.IRI = fmt.Sprintf("process:%s", globalID)
	return &output, nil
}

func mkfifo(procDir string, name string) error {
	path := filepath.Join(procDir, name)
	if err := syscall.Mkfifo(path, 0600); err != nil {
		return fmt.Errorf("making %s fifo: %w", name, err)
	}
	return nil
}

// Taken from the Go internals:
// addCriticalEnv adds any critical environment variables that are required
// (or at least almost always required) on the operating system.
// Currently this is only used for Windows.
func addCriticalEnv(env []string) []string {
	if runtime.GOOS != "windows" {
		return env
	}
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			continue
		}
		k := kv[:eq]
		if strings.EqualFold(k, "SYSTEMROOT") {
			// We already have it.
			return env
		}
	}
	return append(env, "SYSTEMROOT="+os.Getenv("SYSTEMROOT"))
}
