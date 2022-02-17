package shellutil

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/osutil"
)

// Gets the current user's default interactive/login shell environment.
func GetUserEnvironment(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	shellPath := GetUserShellPath()
	switch IdentifyShell(shellPath) {
	case "sh", "bash", "zsh", "fish":
		// All of these are POSIX-compatible and should support the necessary
		// command line flags.
	default:
		// Fallback to POSIX-required standard shell.
		// TODO: Handle Windows shells.
		shellPath = "/bin/sh"
	}

	// Use printenv within boundary lines, in case any gunk is printed during
	// shell initialization.
	boundary := gensym.RandomBase32()
	begin := "BEGIN-" + boundary
	end := "END-" + boundary
	subcmd := fmt.Sprintf("echo %s; printenv; echo %s", begin, end)

	cmd := exec.CommandContext(ctx, shellPath,
		"-ilc", // i=interactive, l=login, c=command.
		subcmd, // Argument to 'c'.
	)
	bs, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewBuffer(bs))
	inside := false
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case begin:
			inside = true
		case end:
			inside = false
		default:
			if inside {
				// NOTE: printenv uses envv, which will produce incorrect results in the
				// presence of newlines or other special characters.
				// TODO: Use an alternative utility that uses DotEnv format.
				name, value := osutil.ParseEnvvEntry(line)
				env[name] = value
			}
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	// Do not include "printenv" or whatever the effective command name is.
	delete(env, "_")

	return env, nil
}
