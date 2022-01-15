package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/deref/exo/internal/about"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionFlags.Verbose, "verbose", "v", false, "Prints version information for all dependencies too.")
}

var versionFlags struct {
	Verbose bool
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of exo",
	Long:  "Print the version of exo.",
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		offline = true
		return cmd.Parent().PersistentPreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			panic("debug.ReadBuildInfo() failed")
		}
		if versionFlags.Verbose {
			fmt.Println(buildInfo.Main.Path, about.Version)
		} else {
			fmt.Println(about.Version)
		}
		if versionFlags.Verbose {
			for _, dep := range buildInfo.Deps {
				fmt.Println(dep.Path, dep.Version)
			}
		}
	},
}
