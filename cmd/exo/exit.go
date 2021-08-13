package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deref/exo/internal/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo daemon",
	Long:  `Stop the exo daemon process.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		loadRunState()
		if runState.Pid == 0 {
			return nil
		}

		return killExod()
	},
}

func killExod() error {
	_ = osutil.TerminateProcessWithTimeout(runState.Pid, 5*time.Second)

	if err := os.Remove(cfg.RunStateFile); err != nil {
		return fmt.Errorf("removing run state file: %w", err)
	}
	return nil
}
