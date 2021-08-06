package main

import (
	"encoding/json"
	"os"

	"github.com/deref/exo/internal/supervise"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(superviseCmd)
}

var superviseCmd = &cobra.Command{
	Hidden: true,
	Use:    "supervise",
	Short:  "Supervises a command",
	Long: `Executes a command, supervises its execution, and redirects stdout/stderr to syslog.

This is an internal use command. See the supervise package implementation for usage details.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &supervise.Config{}
		if err := json.NewDecoder(os.Stdin).Decode(cfg); err != nil {
			cmdutil.Fatalf("reading config from stdin: %v", err)
		}
		if err := cfg.Validate(); err != nil {
			cmdutil.Fatalf("validating config: %v", err)
		}

		supervise.Main(cfg)
	},
}
