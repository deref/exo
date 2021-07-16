package main

import (
	"os"
	"path/filepath"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/osutil"
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
		varDir := cmdutil.MustVarDir()
		pidPath := filepath.Join(varDir, "exod.pid")
		pid := osutil.ReadPid(pidPath)
		if pid == 0 {
			return nil
		}
		process, err := os.FindProcess(pid)
		if err != nil {
			panic(err)
		}
		_ = process.Kill()
		// TODO: Wait for process to exit.
		if err := os.Remove(pidPath); err != nil {
			cmdutil.Fatalf("removing pid file: %w", err)
		}
		return nil
	},
}
