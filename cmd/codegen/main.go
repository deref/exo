package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/deref/exo/josh/codegen"
	"github.com/deref/exo/util/cmdutil"
)

const extension = ".josh.hcl"

func main() {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(d.Name(), extension) {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), extension)

		apiDir := filepath.Dir(path)
		if filepath.Base(apiDir) != "api" {
			return fmt.Errorf("expected %q to be in an api package", path)
		}

		clientDir := filepath.Join(filepath.Dir(apiDir), "client")
		if err := os.Mkdir(clientDir, 0700); err != nil && !os.IsExist(err) {
			return err
		}

		unit, err := codegen.ParseFile(path)
		if err != nil {
			return fmt.Errorf("parsing %q: %w", path, err)
		}

		pkg := &codegen.Package{
			Path: filepath.Join("exo", filepath.Dir(apiDir)),
			Unit: *unit,
		}

		generate := func(dir string, f func(*codegen.Package) ([]byte, error)) error {
			bs, err := f(pkg)
			if err != nil {
				return err
			}
			outpath := filepath.Join(dir, name+".go")
			return ioutil.WriteFile(outpath, bs, 0600)
		}

		if err := generate(apiDir, codegen.GenerateAPI); err != nil {
			return fmt.Errorf("generating %s api: %w", name, err)
		}
		if err := generate(clientDir, codegen.GenerateClient); err != nil {
			return fmt.Errorf("generating %s client: %w", name, err)
		}
		return nil
	})
	if err != nil {
		cmdutil.Fatalf("%v", err)
	}
}
