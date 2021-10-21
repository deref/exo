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

	var converter interface {
		Convert(bs []byte) (*hcl.File, hcl.Diagnostics)
	}
	switch format {
	case "procfile":
		converter = &procfile.Converter{}
	case "compose":
		converter = &compose.Converter{}
	case "exo":
		// No converter needed.
	default:
		return nil, fmt.Errorf("unknown manifest format: %q", l.Format)
	}
	var m *exohcl.Manifest
	if converter == nil {
		m = exohcl.Parse(l.Filename, l.Bytes)
	} else {
		file, diags := converter.Convert(l.Bytes)
		m = exohcl.NewManifest(l.Filename, file, diags)
	}
	var err error
	diags := m.Diagnostics()
	if len(diags) > 0 {
		// Note that this effectively treats all warnings as errors.
		err = diags
	}
	return m, err
}
