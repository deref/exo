package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readmeCmd)
}

var readmeHeading = "## Running exo"

var readmeString = "\n" + readmeHeading + "\nRun the following command to get started:\n```bash\nexo run\n```\n\n"

var readmeCmd = &cobra.Command{
	Use:   "readme [README.md]",
	Short: "Updates your README with instructions on how to run exo",
	Long:  "Updates your README with instructions on how to run exo",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		readmeFile := "README.md"
		if len(args) > 0 {
			readmeFile = args[0]
		}
		byteContents, err := os.ReadFile(readmeFile)
		if err != nil {
			return err
		}
		contents := string(byteContents)
		lines := strings.Split(contents, "\n")
		firstH2 := len(lines) - 1
		for i, line := range lines {
			line := strings.TrimSpace(line)
			if line == readmeHeading {
				return fmt.Errorf("%q already has exo readme", readmeFile)
			}
			if strings.HasPrefix(line, "##") {
				fmt.Println("prefix")
				firstH2 = i
				break
			}
			dashCount := strings.Count(line, "-")
			if i > 1 && dashCount > 0 && dashCount == len(line) {
				fmt.Printf("dashCount: %+v\n", dashCount)
				firstH2 = i - 1
				break
			}
		}

		newContents := strings.Join(lines[:firstH2], "\n") + readmeString + strings.Join(lines[firstH2:], "\n")
		return os.WriteFile(readmeFile, []byte(newContents), 0600)
	},
}
