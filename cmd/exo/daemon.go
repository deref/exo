package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/osutil"
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
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDaemon()
		return nil
	},
}

var runState struct {
	Pid int    `json:"pid"`
	URL string `json:"url"`
}

// checkOrEnsureServer is intended to be used by the CLI commands that do not explicitly start a server.
// It will first check whether a server is running, and if it is not, the behaviour depends on whether
// a URL is configured for the exo client or not:
// - If a URL is configured, this fails since the server is not available.
// - If no URL is configured, this attempts to start a server and exits with an error if it cannot.
// The assumption is that if a URL is configured, the client should not try to start a local instance. One
// implication of this is that if you are running a local server on a non-default port, you will need to
// explicitly run one of the following commands to explicitly start a server: `exo daemon`, `exo server`,
// `exo run`.
func checkOrEnsureServer() {
	if checkHealthy() {
		return
	}
	if cfg.Client.URL == "" {
		ensureDaemon()
		return
	}
	cmdutil.Fatalf("the server at %q is not healthy", cfg.Client.URL)
}

func ensureDaemon() {
	// Validate exod process record.
	err := loadRunState()
	running := false
	switch {
	case err == nil:
		running = osutil.IsValidPid(runState.Pid)
		if !running {
			_ = os.Remove(cfg.RunStateFile)
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
	runState.URL = fmt.Sprintf("http://%s", cmdutil.GetAddr(cfg))
	if err := jsonutil.MarshalFile(cfg.RunStateFile, runState); err != nil {
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
		osutil.TerminateProcessWithTimeout(runState.Pid, time.Second*5)
		os.Exit(1)
	}
}

func loadRunState() error {
	return jsonutil.UnmarshalFile(cfg.RunStateFile, &runState)
}

func newClient() *client.Root {
	url := cfg.Client.URL
	if url == "" {
		url = runState.URL
	}
	// Old state files may contain a url ending in "/".
	url = strings.TrimSuffix(url, "/") + "/_exo/"

	return &client.Root{
		HTTP: http.DefaultClient,
		URL:  url,
	}
}
