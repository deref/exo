package which

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/util/osutil"
)

type Query struct {
	WorkingDirectory string
	PathVariable     string
	Program          string
}

// TODO: Handle PATHEXT on windows for inferring `.exe` extensions, etc.
func (q Query) Run() (string, error) {
	if q.Program == "" {
		return "", errors.New("program is required")
	}
	if filepath.IsAbs(q.Program) {
		if exists, _ := osutil.Exists(q.Program); exists {
			return q.Program, nil
		}
	} else if strings.Contains(q.Program, string(filepath.Separator)) {
		// Relative path.
		candidate := filepath.Join(q.WorkingDirectory, q.Program)
		if exists, _ := osutil.Exists(candidate); exists {
			return candidate, nil
		}
	} else {
		// Search.
		for _, searchPath := range strings.Split(q.PathVariable, string(os.PathListSeparator)) {
			candidate := filepath.Join(searchPath, q.Program)
			if exists, _ := osutil.Exists(candidate); exists {
				return candidate, nil
			}
		}
	}
	return "", fmt.Errorf("%q not found", q.Program)
}

// Uses the ambiant working directory and PATH variable.
func Which(program string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}
	pathVar, _ := os.LookupEnv("PATH")

	return Query{
		WorkingDirectory: wd,
		PathVariable:     pathVar,
		Program:          program,
	}.Run()
}
