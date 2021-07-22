package main

import (
	"fmt"
	"io/ioutil"

	"github.com/deref/exo/exod/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVar(&applyFlags.Format, "format", "", "exo, compose, procfile")
}

var applyFlags struct {
	Format string
}

var applyCmd = &cobra.Command{
	Use:   "apply [flags] [config-file]",
	Short: "Applies a config in the current workspace",
	Long: `Applies a config in the current workspace.

If no config file is specified, a search is conducted in the current directory
in the following order of format preference:

  1. exo
	2. compose
	3. procfile
	
The default exo filename is 'exo.hcl'.

Docker compose files may have one of the following names in order of preference:

	compose.yaml
	compose.yml
  docker-compose.yaml
  docker-compose.yml
	
The expected procfile name 'Procfile'.

If a config format will be guessed from the config filename.  This can be
overidden explicitly with the --format flag.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()

		cl := newClient()
		workspace := requireWorkspace(ctx, cl)

		input := &api.ApplyInput{}
		if len(args) > 0 {
			configPath := args[0]
			input.ConfigPath = &configPath

			// We're not necessarily in the workspace root here,
			// so send the file contents too.
			bs, err := ioutil.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("reading config file: %w", err)
			}
			s := string(bs)
			input.Config = &s
		}
		if applyFlags.Format != "" {
			input.Format = &applyFlags.Format
		}

		_, err := workspace.Apply(ctx, input)
		return err
	},
}
