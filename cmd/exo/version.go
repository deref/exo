package main

import (
	"fmt"
	"runtime/debug"

	"github.com/deref/exo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&versionFlags.Verbose, "verbose", false, "Prints version information for all dependencies too.")
}

var versionFlags struct {
	Verbose bool
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of exo",
	Long:  "Print the version of exo.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			panic("debug.ReadBuildInfo() failed")
		}
		printInfo := func(mod debug.Module) {
			if versionFlags.Verbose {
				fmt.Println(mod.Path, exo.Version)
			} else {
				fmt.Println(exo.Version)
			}
		}
		printInfo(buildInfo.Main)
		if versionFlags.Verbose {
			for _, dep := range buildInfo.Deps {
				fmt.Println(dep.Path, dep.Version)
			}
		}
	},
}
