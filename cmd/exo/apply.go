package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/kernel/api"
	"github.com/deref/exo/osutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVar(&applyFlags.Format, "format", "", "exo, procfile, or compose")
}

var applyFlags struct {
	Format string
}

var applyCmd = &cobra.Command{
	Use:   "apply [flags] [config-file]",
	Short: "Applies a config to the current project",
	Long: `Applies a config to the current project.

If no config file is provided, a search is conducted in the current directory
in the following order of preference:

  exo.hcl
  docker-compose.yml
  Procfile

If a config format will be guessed from the confif filename.  This can be
overidden explicitly with the --format flag.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDeamon()

		configPath := ""
		if len(args) > 0 {
			configPath = args[0]
		}

		if configPath == "" {
			// Search for config.
			for _, candidate := range []string{
				"exo.hcl",
				"docker-compose.yml",
				"Procfile",
			} {
				exist, err := osutil.Exists(candidate)
				if err != nil {
					return fmt.Errorf("searching for config: %w", err)
				}
				if exist {
					configPath = candidate
					break
				}
			}
			if configPath == "" {
				return fmt.Errorf("could not find config file")
			}
		}

		if applyFlags.Format == "" {
			// Guess format.
			name := strings.ToLower(filepath.Base(configPath))
			switch {
			case name == "procfile" || strings.HasSuffix(name, ".procfile"):
				applyFlags.Format = "procfile"
			case name == "docker-compose.yml" || name == "docker-compose.yaml":
				applyFlags.Format = "compose"
			case name == "exo.hcl":
				applyFlags.Format = "exo"
			default:
				return fmt.Errorf("cannot determine config format from name; try the --format flag")
			}
		}

		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("reading config: %w", err)
		}

		client := newClient()

		switch applyFlags.Format {
		case "procfile":
			_, err = client.ApplyProcfile(ctx, &api.ApplyProcfileInput{
				Procfile: string(bs),
			})
		case "compose":
			cmdutil.Fatalf("docker compose configs not yet implemented")
		case "exo":
			_, err = client.Apply(ctx, &api.ApplyInput{
				Config: string(bs),
			})
		}

		return err
	},
}
