package term

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/alessio/shellescape"
)

func EditString(tempPattern string, oldValue string) (string, error) {
	tmpfile, err := ioutil.TempFile("", tempPattern)
	if err != nil {
		return "", fmt.Errorf("creating temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(oldValue); err != nil {
		return "", fmt.Errorf("writing to temporary file: %w", err)
	}
	if err := tmpfile.Close(); err != nil {
		return "", fmt.Errorf("closing temporary file: %w", err)
	}

	editor, err := GetEditor()
	if err != nil {
		return "", fmt.Errorf("getting editor: %w", err)
	}
	if editor == "" {
		return "", errors.New("no editor found")
	}

	// In an attempt to match the behavior of other popular CLI tools, this uses
	// bash, allows passing arguments, handles filenames with spaces and other
	// special shell characters, and bypasses aliases via `command`.
	command := fmt.Sprintf("command %s %s", editor, shellescape.Quote(tmpfile.Name()))
	edit := exec.Command("bash", "-c", command)
	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr
	if err := edit.Run(); err != nil {
		return "", fmt.Errorf("editing file: %w", err)
	}

	newValue, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("reading updated file: %w", err)
	}

	return string(newValue), nil
}

// Returns the preferred command to execute for editing a file.
func GetEditor() (string, error) {
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}
	return LookEditor()
}

// Try to find an editor command in $PATH.
func LookEditor() (string, error) {
	type Candidate struct {
		Name string
		Rest string
	}
	for _, candidate := range []Candidate{
		{"sensible-editor", ""},
		{"code", "--wait"},
		{"nvim", ""},
		{"vim", ""},
		{"nano", ""},
		{"vi", ""},
		{"emacs", ""},
		{"pico", ""},
		{"qe", ""},
		{"mg", ""},
		{"jed", ""},
		{"gedit", ""},
		{"gvim", ""},
		{"ee", ""},
	} {
		found, err := exec.LookPath(candidate.Name)
		if errors.Is(err, exec.ErrNotFound) {
			err = nil
		}
		if err != nil {
			return "", fmt.Errorf("looking for %q: %w", candidate, err)
		}
		if found != "" {
			return fmt.Sprintf("%s %s", found, candidate.Rest), nil
		}
	}
	return "", nil
}
