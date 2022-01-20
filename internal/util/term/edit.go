package term

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/deref/exo/internal/util/which"
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

	stat, err := os.Stat(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("checking modification time: %w", err)
	}
	originalModTime := stat.ModTime()

	editor := os.Getenv("EDITOR")
	if editor == "" {

		for _, candidateEditor := range []string{
			"sensible-editor",
			"editor",
			"code",
			"vim",
			"nano",
			"vi",
			"emacs",
			"ee",
		} {
			found, err := which.Which(candidateEditor)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				return "", fmt.Errorf("looking up candidate editor %q: %w", candidateEditor, err)
			}
			if found != "" {
				editor = found
				break
			}
		}

		if editor == "" {
			return "", errors.New("no editor available")
		}
	}

	edit := exec.Command(editor, tmpfile.Name())
	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr
	if err := edit.Run(); err != nil {
		return "", fmt.Errorf("editing file: %w", err)
	}

	if stat, err = os.Stat(tmpfile.Name()); err != nil {
		return "", fmt.Errorf("checking modification time: %w", err)
	}
	newModTime := stat.ModTime()
	if !newModTime.After(originalModTime) {
		fmt.Fprintf(os.Stderr, "file unmodified - not updating.\n")
		return oldValue, nil
	}

	newValue, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("reading updated file: %w", err)
	}

	return string(newValue), nil
}
