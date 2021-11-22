package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/osutil"
)

type Resource struct {
	Pid int
}

type Start struct {
	// Absolute path to working directory.
	Directory string `json:"directory"`
	// Absolute path to program to execute.
	Program string `json:"program"`
	// Command line arguments.
	Arguments []string `json:"arguments"`
	// Complete environment given to the process.
	Environment map[string]string `json:"environment"`
	// File paths to attach. First three are generally stdin, stdout, and stderr.
	Files []OpenFile `json:"files"`
}

type OpenFile struct {
	Path string `json:"path"`
	Flag int    `json:"flag"`
	Perm int    `json:"perm"`
}

func (res *Resource) Create(ctx context.Context, spec string) (iri string, err error) {
	var start Start
	if err := jsonutil.UnmarshalString(spec, &start); err != nil {
		return "", fmt.Errorf("unmarshalling spec: %w", err)
	}

	argv := append([]string{start.Program}, start.Arguments...)

	envKeys := make([]string, 0, len(start.Environment))
	for k := range start.Environment {
		envKeys = append(envKeys, k)
	}
	sort.Strings(envKeys)
	env := make([]string, len(start.Environment))
	for i, k := range envKeys {
		env[i] = fmt.Sprintf("%s=%s", k, start.Environment[k])
	}

	files := make([]*os.File, len(start.Files))
	for fd, startFile := range start.Files {
		f, err := os.OpenFile(startFile.Path, startFile.Flag, os.FileMode(startFile.Perm))
		if err != nil {
			return "", fmt.Errorf("opening fd %d: %w", fd, err)
		}
		files[fd] = f
	}

	proc, err := os.StartProcess(argv[0], argv, &os.ProcAttr{
		Dir:   start.Directory,
		Env:   env,
		Files: files,
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	})
	if err != nil {
		return "", err
	}

	pid := proc.Pid
	iri = MakeIRI(pid)
	return iri, nil
}

func (res *Resource) LoadIRI(ctx context.Context, iri string) error {
	pid, err := ParseIRI(iri)
	if err != nil {
		return fmt.Errorf("invalid process iri: %w", err)
	}
	res.Pid = pid
	return nil
}

func MakeIRI(pid int) string {
	return fmt.Sprintf("exo:unix:process/%d", pid) // XXX uri scheme ok here?
}

func ParseIRI(iri string) (pid int, err error) {
	parts := strings.SplitN(iri, "/", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected /")
	}
	pid64, err := strconv.ParseInt(parts[1], 10, 32)
	return int(pid64), err
}

func (res *Resource) Exists(ctx context.Context) (bool, error) {
	return osutil.IsValidPid(res.Pid), nil
}

func (res *Resource) Update(ctx context.Context, patch string) error {
	return errors.New("unsupported patch")
}

func (res *Resource) Delete(ctx context.Context) error {
	if err := osutil.SignalGroup(res.Pid, os.Interrupt); err != nil {
		return fmt.Errorf("interrupting: %w", err)
	}

	ok := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = osutil.KillGroup(res.Pid)
		case <-ok:
		}
	}()

	if _, err := osutil.WaitProcess(res.Pid); err != nil {
		return fmt.Errorf("waiting: %w", err)
	}
	close(ok)
	return nil
}
