package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVar(&convertFlags.Format, "format", "", "compose, procfile")
}

var convertFlags struct {
	Format string
}

var convertCmd = &cobra.Command{
	Use:   "convert [flags] [manifest-file]",
	Short: "Converts a docker-compose file or Procfile into an exo manifest",
	Long: `Converts a docker-compose file or Procfile into an exo manifest

	If unspecified, a manifest format will be guessed from the manifest filename.  This can be
	overidden explicitly with the --format flag.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath := args[0]
		if !path.IsAbs(manifestPath) {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting working dir: %w", err)
			}
			manifestPath = path.Join(wd, manifestPath)
		}

		format := &convertFlags.Format
		if convertFlags.Format == "" {
			format = nil
		}
		filename, err := convertAndSave(manifestPath, format)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Converted manifest and saved to %s\n", filename)
		return nil
	},
}

func convertAndSave(manifestPath string, manifestFormat *string) (string, error) {
	if path.Base(manifestPath) == "exo.json" || path.Base(manifestPath) == "exo.yaml" || (manifestFormat != nil && *manifestFormat == "exo") {
		return manifestPath, nil
	}

	res := loadManifest(manifestPath, manifestFormat)
	if res.Err != nil {
		return "", res.Err
	}

	for _, warn := range res.Warnings {
		fmt.Printf("WARNING: %s\n", warn)
	}

	manifestBytes, err := yaml.MarshalWithOptions(
		res.Manifest,
		yaml.IndentSequence(true),
	)
	if err != nil {
		return "", fmt.Errorf("marshalling manifest: %w", err)
	}

	dir, _ := path.Split(manifestPath)
	filename := path.Join(dir, "exo.yaml")
	return filename, ioutil.WriteFile(filename, manifestBytes, 0600)
}

type manifestCandidate struct {
	Format   string
	Filename string
}

var manifestCandidates = []manifestCandidate{
	{"compose", "compose.yaml"},
	{"compose", "compose.yml"},
	{"compose", "docker-compose.yaml"},
	{"compose", "docker-compose.yml"},
	{"procfile", "Procfile"},
}

func loadManifest(manifestPath string, manifestFormat *string) manifest.LoadResult {
	manifestString := ""

	bs, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return manifest.LoadResult{
			Err: fmt.Errorf("reading manifest file: %w", err),
		}
	}
	manifestString = string(bs)

	format := ""
	if manifestFormat == nil {
		// Guess format.
		name := strings.ToLower(filepath.Base(manifestPath))
		switch name {
		case "procfile":
			format = "procfile"
		case "compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml":
			format = "compose"
		default:
			if strings.HasPrefix(name, "procfile.") || strings.HasSuffix(name, ".procfile") {
				format = "procfile"
			} else {
				return manifest.LoadResult{
					Err: fmt.Errorf("cannot determine manifest format from file name: %s", name),
				}
			}
		}
	} else {
		format = *manifestFormat
	}

	var loader interface {
		Load(r io.Reader) manifest.LoadResult
	}
	switch format {
	case "procfile":
		loader = procfile.Loader
	case "compose":
		dir, _ := path.Split(manifestPath)
		projectName := path.Base(dir)
		projectName = manifest.MangleName(projectName)
		loader = &compose.Loader{ProjectName: projectName}
	case "exo":
		loader = manifest.Loader
	default:
		return manifest.LoadResult{
			Err: fmt.Errorf("unknown manifest format: %q", format),
		}
	}

	res := loader.Load(strings.NewReader(manifestString))
	// TODO: Validate manifest.
	return res
}
