package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/osutil"
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
		ensureDeamon()

		configPath := ""
		if len(args) > 0 {
			configPath = args[0]
		}

		if configPath == "" {
			// Search for config.
			for _, candidate := range []string{
				"exo.hcl",
				"compose.yaml",
				"compose.yml",
				"docker-compose.yaml",
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
			switch name {
			case "procfile":
				applyFlags.Format = "procfile"
			case "compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml":
				applyFlags.Format = "compose"
			case "exo.hcl":
				applyFlags.Format = "exo"
			default:
				if strings.HasSuffix(name, ".procfile") {
					applyFlags.Format = "procfile"
				} else {
					return fmt.Errorf("cannot determine config format from name; try the --format flag")
				}
			}
		}

		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("reading config: %w", err)
		}

		cl := newClient()
		workspace := requireWorkspace(ctx, cl)

		switch applyFlags.Format {
		case "procfile":
			_, err = workspace.ApplyProcfile(ctx, &api.ApplyProcfileInput{
				Procfile: string(bs),
			})
		case "compose":
			cmdutil.Fatalf("docker compose configs not yet implemented")
		case "exo":
			_, err = workspace.Apply(ctx, &api.ApplyInput{
				Config: string(bs),
			})
		}

		return err
	},
}
