package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/deref/exo/internal/resolvers"
	"github.com/deref/exo/internal/util/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Report exo status",
	Long:  `Report exo daemon status.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := loadRunState()
		if err != nil {
			return fmt.Errorf("checking run state: %w", err)
		}

		pid := ""
		if osutil.IsValidPid(runState.Pid) {
			pid = strconv.Itoa(runState.Pid)
		}
		healthy, _ := checkHealthAndVersion()

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "healthy:\t%t\n", healthy)
		_, _ = fmt.Fprintf(w, "pid:\t%s\n", pid)
		_, _ = fmt.Fprintf(w, "gui:\t%s\n", effectiveServerURL())
		_ = w.Flush()
		return nil
	},
}

func checkHealthAndVersion() (bool, string) {
	res, err := http.Get(effectiveServerURL() + "/_exo/health")
	if err != nil {
		return false, ""
	}
	bs, _ := ioutil.ReadAll(res.Body)
	version := res.Header.Get("Exo-Version")
	// See note [HEALTH_CHECK].
	return string(bytes.TrimSpace(bs)) == "ok", version
}

// XXX This is confused for API and GUI, which may not be on the same port.
func effectiveServerURL() string {
	url := cfg.Client.URL
	if url == "" {
		url = runState.URL
	}
	if url == "" {
		url = fmt.Sprintf("http://localhost:%d", cfg.GUI.Port)
	}
	return url
}

// XXX Do not utilize resolvers directly.
func newGUIRoutes() *resolvers.RoutesResolver {
	return resolvers.NewRoutesResolver(effectiveServerURL())
}
