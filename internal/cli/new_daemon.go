package cli

import (
	"strings"

	"github.com/deref/exo/internal/providers/os/components/daemon"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.AddCommand(newDaemonCmd)
	newDaemonCmd.Flags().StringVarP(
		&daemonSpec.Directory,
		"directory", "d",
		"",
		"set the working directory for the daemon",
	)
}

var daemonSpec = daemon.Spec{}

var newDaemonCmd = &cobra.Command{
	Use:   "daemon <name> [options] [--] [name=value ...] <program> [args ...]",
	Short: "Creates a new daemon",
	Long: `Creates a new daemon component, running a host os process.

The double dash separator is recommended to avoid flag confusion between
exo flags and options for your program.

Environment variables may be specified by providing name=value pairs
before the program name.
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		name := args[0]
		args = args[1:]

		for len(args) > 0 {
			arg := args[0]
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) < 2 {
				break
			}
			args = args[1:]
			name := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if daemonSpec.Environment == nil {
				daemonSpec.Environment = make(map[string]string)
			}
			daemonSpec.Environment[name] = value
		}

		if len(args) == 0 {
			return cmd.Usage()
		}
		daemonSpec.Program = args[0]
		daemonSpec.Arguments = args[1:]

		return createComponent(ctx, name, "daemon", daemonSpec)
	},
}
