package cmdutil

import "os"

func NonTerminalStdio() bool {
	for _, device := range []*os.File{os.Stdin, os.Stdout, os.Stderr} {
		if fileInfo, _ := device.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}
