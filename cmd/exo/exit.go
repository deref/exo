package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo daemon",
	Long:  `Stop the exo daemon process.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		loadRunState()
		if runState.Pid == 0 {
			return nil
		}

		return killExod()
	},
}

func killExod() error {
	process, err := os.FindProcess(runState.Pid)
	if err != nil {
		panic(err)
	}
	// TODO: Try to stop gracefully.
	_ = process.Kill()

	// TODO: Wait for process to exit.

	if err := os.Remove(knownPaths.RunStateFile); err != nil {
		return fmt.Errorf("removing run state file: %w", err)
	}
	return nil
}
