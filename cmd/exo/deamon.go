package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/api"
	"github.com/deref/exo/kernel/client"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/osutil"
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
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return nil
	},
}

var runState struct {
	Pid int    `json:"pid"`
	URL string `json:"url"`
}

func ensureDeamon() {
	paths := cmdutil.MustMakeDirectories()

	// Validate exod process record.
	err := loadRunState(paths.RunStateFile)
	running := false
	switch {
	case err == nil:
		running = osutil.IsValidPid(runState.Pid)
		if !running {
			_ = os.Remove(paths.RunStateFile)
		}
	case os.IsNotExist(err):
		// Not running.
	default:
		cmdutil.Fatalf("checking run state: %w", err)
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

	// Write run state.
	runState.Pid = cmd.Process.Pid
	runState.URL = fmt.Sprintf("http://%s/", cmdutil.GetAddr())
	if err := jsonutil.MarshalFile(paths.RunStateFile, runState); err != nil {
		cmdutil.Fatalf("writing run state: %w", err)
	}
}

func loadRunState(path string) error {
	return jsonutil.UnmarshalFile(path, &runState)
}

func newClient() api.Kernel {
	return client.NewProject(&josh.Client{
		HTTP: http.DefaultClient,
		URL:  runState.URL + "_exo/",
	})
}
