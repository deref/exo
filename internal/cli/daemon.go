package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/deref/exo/internal/about"
	gqlclient "github.com/deref/exo/internal/client"
	"github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/exod"
	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/jmoiron/sqlx"
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
		if cfg.NoDaemon {
			cmdutil.Fatalf("daemon disabled by config")
		}
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
	healthy, version := checkHealthAndVersion()
	outOfDate := healthy && (version != about.Version)
	if outOfDate {
		cmdutil.Fatalf("daemon at %q is not up to date. Please restart the server.\ndaemon version: %s\nclient version: %s", effectiveServerURL(), version, about.Version)
	}
	if healthy {
		return
	}
	if cfg.NoDaemon {
		cmdutil.Fatalf("daemon at %q is not healthy\ndeamonization disabled by config", effectiveServerURL())
	}
	if cfg.Client.URL == "" {
		ensureDaemon()
		return
	}
	cmdutil.Fatalf("the server at %q is not healthy", cfg.Client.URL)
}

func ensureDaemon() {
	exod.MustMakeDirectories(cfg)

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
	delay := 10 * time.Millisecond
	for attempt := 0; attempt < 15; attempt++ {
		if healthy, _ := checkHealthAndVersion(); healthy {
			ok = true
			break
		}
		<-time.After(delay)
		if delay < 100 {
			delay *= 2
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

func mustGetToken() string {
	token, err := cfg.GetTokenClient().GetToken()
	if err != nil {
		cmdutil.Fatalf("getting token client: %w", err)
	}
	return token
}

func clientURL() string {
	url := cfg.Client.URL
	if url == "" {
		url = runState.URL
	}
	// Old state files may contain a url ending in "/".
	return strings.TrimSuffix(url, "/")
}

// TODO: Remove me after switching fully to graphql.
func newClient() *client.Root {
	return &client.Root{
		HTTP:  http.DefaultClient,
		URL:   clientURL() + "/_exo/",
		Token: mustGetToken(),
	}
}

func dialGraphQL(ctx context.Context) (client *gqlclient.Client, shutdown func()) {
	// XXX this is a hack for testing daemonless. See exod/main.go & reconcile with that.
	dbPath := filepath.Join(cfg.VarDir, "exo.sqlite3")
	txMode := "exclusive"
	connStr := dbPath + "?_txlock=" + txMode
	db, err := sqlx.Open("sqlite3", connStr)
	if err != nil {
		cmdutil.Fatalf("opening sqlite db: %v", err)
	}
	shutdown = func() {
		if err := db.Close(); err != nil {
			cmdutil.Warnf("error closing sqlite db: %v", err)
		}
	}

	root := &resolvers.RootResolver{
		DB: db,
	}
	if err := root.Migrate(ctx); err != nil {
		cmdutil.Fatalf("migrating db: %w", err)
	}

	httpClient := &http.Client{
		Transport: &httputil.NetworklessTransport{
			Handler: resolvers.NewHandler(root),
		},
	}
	client = gqlclient.NewClient(clientURL()+"/graphql/", httpClient)
	return
}
