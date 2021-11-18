package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/util/osutil"
	"github.com/deref/exo/internal/util/which"
	"github.com/spf13/cobra"
)

func init() {
	completionCmd.AddCommand(completionInstallCmd)
}

var completionInstallCmd = &cobra.Command{
	Use:   "install [bash|zsh|fish|powershell]",
	Short: "Install shell completions",
	Args:  cobra.MaximumNArgs(1),
	Long: `Installs shell completions.

If the shell argument is not provided, it will be inferred from the SHELL
environment variable.

After running this command, you must restart your shell. Some shells, such as
zsh, may require additional steps to clear completion caches.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := ""
		if len(args) > 0 {
			shell = args[0]
		}

		if shell == "" {
			shell = inferShell()
		}
		if shell == "" {
			return errors.New("cannot determine shell")
		}
		fmt.Println("inferred shell:", shell)

		var completionPath string
		for _, candidate := range completionPathCandidates(shell) {
			dir := filepath.Dir(candidate)
			if writable, _ := osutil.IsWritable(dir); writable {
				completionPath = candidate
				break
			}
		}

		if completionPath == "" {
			return errors.New("cannot infer completion installation path")
		}
		f, err := os.Create(completionPath)
		if err != nil {
			return fmt.Errorf("creating completion file: %w", err)
		}
		defer f.Close()
		completionGenerate(f, shell)

		fmt.Println("wrote completion script:", completionPath)
		return nil
	},
}

func completionPathCandidates(shell string) []string {
	var paths []string
	switch shell {
	case "bash":
		paths = []string{
			completionPathBashLinux,
			completionPathBashMac,
		}
	case "zsh":
		paths = []string{
			completionPathZsh(),
		}
	case "fish":
		paths = []string{
			completionPathFish(),
		}
	case "powershell":
		// TODO: Figure out how to auto-install powershell.
	}

	// Remove empty string paths.
	dest := 0
	for i := 0; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		paths[dest] = paths[i]
		dest++
	}
	return paths[:dest]
}

// Returns the path of the current user's default shell.
func getUserShell() string {
	sudoUser, _ := os.LookupEnv("SUDO_USER")
	if sudoUser != "" {
		// TODO: Query the user's shell from getent or similar on Linux.
		return ""
	}
	shell, _ := os.LookupEnv("SHELL")
	return shell
}

// Returns the name of the current user's default shell.
func inferShell() string {
	shellPath := getUserShell()
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

const completionPathBashLinux = "/etc/bash_completion.d/exo"
const completionPathBashMac = "/usr/local/etc/bash_completion.d/exo"

func completionPathZsh() string {
	shell := getUserShell()
	if !strings.HasSuffix(shell, "/zsh") {
		shell, _ = which.Which("zsh")
	}
	if shell == "" {
		return ""
	}
	bs, _ := exec.Command(shell, "-c", `echo "${fpath[1]}/_exo"`).Output()
	return string(bytes.TrimSpace(bs))
}

func completionPathFish() string {
	// See https://fishshell.com/docs/current/completions.html
	completionFile := "/usr/share/fish/vendor_completions.d/exo.fish"
	cmd := exec.Command("pkg-config", "--variable", "completionsdir", "fish")
	output, err := cmd.Output()
	if err != nil {
		return completionFile
	}
	return string(bytes.TrimSpace(output))
}
