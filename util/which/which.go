package which

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Query struct {
	WorkingDirectory string
	PathVariable     string
	Program          string
}

func (q *Query) Run() (string, error) {
	program := q.Program
	if program == "" {
		return "", errors.New("program is required")
	}
	if strings.Contains(program, string(filepath.Separator)) {
		// Relative path.
		program = filepath.Join(q.WorkingDirectory, program)
	} else {
		// Search.
		for _, searchPath := range strings.Split(q.PathVariable, string(os.PathListSeparator)) {
			candidate := filepath.Join(searchPath, program)
			info, _ := os.Stat(candidate)
			if info != nil {
				program = candidate
				break
			}
		}
	}
	// TODO: Handle PATHEXT on windows for inferring `.exe` extensions, etc.
	if program == "" {
		return "", fmt.Errorf("%q not found", q.Program)
	}
	return program, nil
}
