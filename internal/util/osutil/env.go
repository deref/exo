package osutil

import (
	"fmt"
	"sort"

	"github.com/alessio/shellescape"
)

// Convert an environment map to a slice of strings suitable for exec.
// Entries are of the form name=value and sorted. No escaping is performed.
func EnvMapToEnvv(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for name, value := range m {
		s[i] = fmt.Sprintf("%s=%s", name, value)
		i++
	}
	sort.Strings(s)
	return s
}

// Convert an environment map to a slice of strings suitable for a .env file.
// Entries are of the form name=value and sorted. Values are shell-quoted.
func EnvMapToDotEnv(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for name, value := range m {
		s[i] = fmt.Sprintf("%s=%s", name, shellescape.Quote(value))
		i++
	}
	sort.Strings(s)
	return s
}
