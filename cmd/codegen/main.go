package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/deref/exo/josh/codegen"
	"github.com/deref/exo/josh/idl"
	"github.com/deref/exo/josh/model"
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

		pkgPath := filepath.Join("exo", filepath.Dir(apiDir))
		pkg := model.NewPackage(pkgPath)
		idl.LoadFile(pkg, path)
		if err := pkg.Err(); err != nil {
			return fmt.Errorf("loading %q: %w", path, err)
		}

		generate := func(dir string, f func(*model.Package) ([]byte, error)) error {
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
