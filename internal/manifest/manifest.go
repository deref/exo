package manifest

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/deref/exo/internal/util/errutil"
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

type Loader struct {
	WorkspaceName string
	Format        string
	Filename      string
	Reader        io.Reader
}

func (l *Loader) Load() (*exohcl.Manifest, error) {
	format := l.Format
	if format == "" {
		if l.Filename == "" || l.Filename == "/dev/stdin" {
			format = "exo"
		} else {
			format = GuessFormat(l.Filename)
			if format == "" {
				return nil, errutil.NewHTTPError(http.StatusBadRequest, "cannot determine manifest format from file name")
			}
		}
	}

	var formatLoader interface {
		Load(r io.Reader) (*exohcl.Manifest, error)
	}
	switch format {
	case "procfile":
		formatLoader = procfile.Loader
	case "compose":
		formatLoader = &compose.Loader{
			ProjectName: l.WorkspaceName,
		}
	case "exo":
		formatLoader = &exohcl.Loader{
			Filename: l.Filename,
		}
	default:
		return nil, fmt.Errorf("unknown manifest format: %q", l.Format)
	}
	return formatLoader.Load(l.Reader)
}
