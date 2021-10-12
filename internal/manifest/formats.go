package manifest

import (
	"path/filepath"
	"strings"
)

func GuessFormat(path string) string {
	name := strings.ToLower(filepath.Base(path))
	switch name {
	case "exo.hcl":
		return "exo"
	case "procfile":
		return "procfile"
	case "compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml":
		return "compose"
	default:
		if strings.HasPrefix(name, "procfile.") || strings.HasSuffix(name, ".procfile") {
			return "procfile"
		}
		return ""
	}
}
