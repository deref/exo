package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/pathutil"
)

func (ws *Workspace) loadManifest(rootDir string, input *api.ApplyInput) manifest.LoadResult {
	manifestString := ""
	manifestPath := path.Join(rootDir, "exo.yaml")
	if input.ManifestPath != nil {
		manifestPath = *input.ManifestPath
	}
	if input.Manifest == nil {
		if !pathutil.HasFilePathPrefix(manifestPath, rootDir) {
			return manifest.LoadResult{
				Err: errors.New("cannot read manifest outside of workspace root"),
			}
		}

		bs, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return manifest.LoadResult{
				Err: fmt.Errorf("reading manifest file: %w", err),
			}
		}
		manifestString = string(bs)
	} else {
		manifestString = *input.Manifest
	}

	res := manifest.Loader.Load(strings.NewReader(manifestString))
	if res.Err != nil {
		res.Err = errutil.WithHTTPStatus(http.StatusBadRequest, res.Err)
	}
	// TODO: Validate manifest.
	return res
}
