package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	completionCmd.AddCommand(completionGenerateCmd)
}

var completionGenerateCmd = &cobra.Command{
	Use:   "generate [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `For automatic installation, see:

	$ exo completion install --help

To load completions:

Bash:

  $ source <(exo completion generate bash)

  # To load completions for each session, execute once:
  # Linux:
  $ exo completion generate bash > /etc/bash_completion.d/exo
  # macOS:
  $ exo completion generate bash > /usr/local/etc/bash_completion.d/exo

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ exo completion generate zsh > "${fpath[1]}/_exo"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ exo completion generate fish | source

  # To load completions for each session, execute once:
  $ exo completion generate fish > ~/.config/fish/completions/exo.fish

PowerShell:

  PS> exo completion generate powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> exo completion generate powershell > exo.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		completionGenerate(os.Stdout, args[0])
	},
}

func completionGenerate(w io.Writer, shell string) {
	switch shell {
	case "bash":
		rootCmd.GenBashCompletion(w)
	case "zsh":
		rootCmd.GenZshCompletion(w)
	case "fish":
		rootCmd.GenFishCompletion(w, true)
	case "powershell":
		rootCmd.GenPowerShellCompletion(w)
	}
}
