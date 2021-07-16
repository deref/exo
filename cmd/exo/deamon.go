package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deamonCmd)
}

var deamonCmd = &cobra.Command{
	Hidden: true,
	Use:    "deamon",
	Short:  "Start the exo deamon",
	Long: `Start the exo deamon and then do nothing else.

Since most commands implicitly start the exo deamon, users generally do not
have to invoke this themselves.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return nil
	},
}

func ensureDeamon() {
	varDir := cmdutil.MustVarDir()

	// Validate pid file.
	pidPath := filepath.Join(varDir, "exod.pid")
	_, err := os.Stat(pidPath)
	running := false
	switch {
	case err == nil:
		pid := osutil.ReadPid(pidPath)
		running = osutil.IsValidPid(pid)
		if !running {
			_ = os.Remove(pidPath)
		}
	case os.IsNotExist(err):
		// Not running.
	default:
		cmdutil.Fatalf("checking pid file: %w", err)
	}

	if running {
		// TODO: health check.
		return
	}

	// Start server in background.
	exoPath := os.Args[0]
	cmd := exec.Command(exoPath, "server")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		cmdutil.Fatalf("starting exo server: %w", err)
	}

	// Write pid file.
	pid := cmd.Process.Pid
	if err := osutil.WritePid(pidPath, pid); err != nil {
		cmdutil.Fatalf("writing pid: %w", err)
	}
}
