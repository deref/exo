package main

import (
	"os"

	"github.com/deref/exo/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo deamon",
	Long:  `Stop the exo deamon process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := cmdutil.MustMakeDirectories()
		loadRunState(paths.RunStateFile)
		if runState.Pid == 0 {
			return nil
		}

		process, err := os.FindProcess(runState.Pid)
		if err != nil {
			panic(err)
		}
		_ = process.Kill()

		// TODO: Wait for process to exit.

		if err := os.Remove(paths.RunStateFile); err != nil {
			cmdutil.Fatalf("removing pid file: %w", err)
		}
		return nil
	},
}
