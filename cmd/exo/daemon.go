package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/deref/exo/exod/client"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/jsonutil"
	"github.com/deref/exo/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Hidden: true,
	Use:    "daemon",
	Short:  "Start the exo daemon",
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
	cmd := exec.Command(exoPath, "server") // TODO: CommandContext.
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

	// Wait server to be healthy.
	ok := false
	var delay time.Duration
	for attempt := 0; attempt < 15; attempt++ {
		<-time.After(delay)
		delay = 20 * time.Millisecond
		res, err := http.Get(runState.URL + "_exo/health")
		if err != nil {
			continue
		}
		bs, _ := ioutil.ReadAll(res.Body)
		// See note [HEALTH_CHECK].
		if string(bytes.TrimSpace(bs)) == "ok" {
			ok = true
			break
		}
	}

	// Cleanup if unhealthy.
	if !ok {
		cmdutil.Warnf("daemon not healthy")
		killExod(paths)
		os.Exit(1)
	}
}

func loadRunState(path string) error {
	return jsonutil.UnmarshalFile(path, &runState)
}

func newClient() *client.Root {
	return &client.Root{
		HTTP: http.DefaultClient,
		URL:  runState.URL + "_exo/",
	}
}
