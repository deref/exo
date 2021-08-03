package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/deref/exo/core/client"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the exo daemon",
	Long: `Start the exo daemon and then do nothing else.

Since most commands implicitly start the exo daemon, users generally do not
have to invoke this themselves.`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDaemon()
		return nil
	},
}

var runState struct {
	Pid int    `json:"pid"`
	URL string `json:"url"`
}

func ensureDaemon() {
	// Validate exod process record.
	err := loadRunState()
	running := false
	switch {
	case err == nil:
		running = osutil.IsValidPid(runState.Pid)
		if !running {
			_ = os.Remove(knownPaths.RunStateFile)
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
	cmd := exec.Command(exoPath, "server") // TODO: CommandContext.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		cmdutil.Fatalf("starting exo server: %w", err)
	}

	// Write run state.
	runState.Pid = cmd.Process.Pid
	runState.URL = fmt.Sprintf("http://%s/", cmdutil.GetAddr(cfg))
	if err := jsonutil.MarshalFile(knownPaths.RunStateFile, runState); err != nil {
		cmdutil.Fatalf("writing run state: %w", err)
	}

	// Wait server to be healthy.
	ok := false
	var delay time.Duration
	for attempt := 0; attempt < 15; attempt++ {
		<-time.After(delay)
		delay = 20 * time.Millisecond
		if checkHealthy() {
			ok = true
			break
		}
	}

	// Cleanup if unhealthy.
	if !ok {
		cmdutil.Warnf("daemon not healthy")
		killExod()
		os.Exit(1)
	}
}

func loadRunState() error {
	return jsonutil.UnmarshalFile(knownPaths.RunStateFile, &runState)
}

func newClient() *client.Root {
	return &client.Root{
		HTTP: http.DefaultClient,
		URL:  runState.URL + "_exo/",
	}
}
