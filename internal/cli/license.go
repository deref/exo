package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/deref/exo/internal/about"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(licenseCmd)
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Display license and legal notices",
	Long:  `Displays license and required legal notices.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		offline = true
		return cmd.Parent().PersistentPreRunE(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := about.Notices.ReadDir(".")
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.Type().IsDir() {
				continue
			}
			name := entry.Name()
			hr := strings.Repeat("-", len(name))
			fmt.Println(hr)
			fmt.Println(entry.Name())
			fmt.Println(hr)
			fmt.Println()
			bs, err := about.Notices.ReadFile(name)
			if err != nil {
				panic(err)
			}
			os.Stdout.Write(bs)
			fmt.Println()
		}
	},
}
