package osutil

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
)

// Convert an environment map to a slice of strings suitable for exec.
// Entries are of the form name=value and sorted. No escaping is performed.
func EnvMapToEnvv(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for name, value := range m {
		s[i] = FormatEnvvEntry(name, value)
		i++
	}
	sort.Strings(s)
	return s
}

// Returns an entry in the form name=value. No escaping is performed.
func FormatEnvvEntry(name, value string) string {
	return fmt.Sprintf("%s=%s", name, value)
}

// Parses an name=value entry where the value is returned literally.
func ParseEnvvEntry(entry string) (name, value string) {
	parts := strings.SplitN(entry, "=", 2)
	return parts[0], parts[1]
}

// Convert an environment map to a slice of strings suitable for a .env file.
// Entries are of the form name=value and sorted. Values are shell-quoted.
func EnvMapToDotEnv(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for name, value := range m {
		s[i] = FormatDotEnvEntry(name, value)
		i++
	}
	sort.Strings(s)
	return s
}

// Creates entries of the form name=value. Values are shell-quoted.
func FormatDotEnvEntry(name, value string) string {
	return fmt.Sprintf("%s=%s", name, shellescape.Quote(value))
}
