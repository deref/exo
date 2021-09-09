package main

import (
	"fmt"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/unix/components/process"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.AddCommand(newProcessCmd)
	newProcessCmd.Flags().StringVarP(
		&processSpec.Directory,
		"directory", "d",
		"",
		"set the working directory for the process",
	)
}

var processSpec = process.Spec{}

var newProcessCmd = &cobra.Command{
	Use:   "process <name> [options] [--] [name=value ...] <program> [args ...]",
	Short: "Creates a new process",
	Long: `Creates a new process.
	
The double dash separator is recommended to avoid flag confusion between
exo flags and options for your program.
	
Environment variables may be specified by providing name=value pairs
before the program name.
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)

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
			if processSpec.Environment == nil {
				processSpec.Environment = make(map[string]string)
			}
			processSpec.Environment[name] = value
		}

		if len(args) == 0 {
			return cmd.Usage()
		}
		processSpec.Program = args[0]
		processSpec.Arguments = args[1:]

		output, err := workspace.CreateComponent(ctx, &api.CreateComponentInput{
			Name: name,
			Type: "process",
			Spec: jsonutil.MustMarshalString(processSpec),
		})
		if err != nil {
			return err
		}
		fmt.Println(output.ID)
		return nil
	},
}
