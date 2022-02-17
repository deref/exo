package shellutil

import (
	"os"
	"strings"
)

// Returns the path of the current user's default shell.
func GetUserShellPath() string {
	sudoUser, _ := os.LookupEnv("SUDO_USER")
	if sudoUser != "" {
		// TODO: Query the user's shell from getent or similar on Linux.
		return ""
	}
	shell, _ := os.LookupEnv("SHELL")
	return shell
}

// Attempts to return the name of the current user's default shell.
func IdentifyUserShell() string {
	return IdentifyShell(GetUserShellPath())
}

// Given a shell path, attempts to returns the name of a known shell
// implementation.
func IdentifyShell(shellPath string) string {
	for _, shellName := range []string{
		"bash",
		"zsh",
		"fish",
	} {
		if strings.HasSuffix(shellPath, "/"+shellName) {
			return shellName
		}
	}
	// TODO: Detect powershell.
	return ""
}
