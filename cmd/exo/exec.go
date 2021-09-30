package main

import (
	"fmt"
	"syscall"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/which"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:                   "exec [flags] -- <program> [argument ...]",
	Short:                 "Execute program with environment",
	Long:                  `Executes program in the current workspace's environment.`,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		envv, err := getEnvv()
		if err != nil {
			return fmt.Errorf("getting environment: %w", err)
		}
		program, err := which.Which(args[0])
		if err != nil {
			return fmt.Errorf("resolving program: %w", err)
		}
		err = syscall.Exec(program, args, envv)
		cmdutil.Fatalf("%v", err)
		panic("unreachable")
	},
}
