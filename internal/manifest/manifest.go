package manifest

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/exohcl"
	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/hashicorp/hcl/v2"
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
	Bytes         []byte
}

func (l *Loader) Load(ctx *exohcl.AnalysisContext) (*exohcl.Manifest, error) {
	m := &exohcl.Manifest{
		Filename: l.Filename,
	}

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

	var importer interface {
		Import(ctx *exohcl.AnalysisContext, bs []byte) *hcl.File
	}
	switch format {
	case "procfile":
		importer = &procfile.Importer{}
	case "compose":
		importer = &compose.Importer{
			ProjectName: l.WorkspaceName,
		}
	case "exo":
		importer = &exohcl.Importer{
			Filename: l.Filename,
		}
	default:
		return m, fmt.Errorf("unknown manifest format: %q", l.Format)
	}
	m.File = importer.Import(ctx, l.Bytes)
	m.Analyze(ctx)
	return m, nil
}
