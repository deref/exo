package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/josh/codegen"
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
		module, err := codegen.ParseFile(path)
		if err != nil {
			return fmt.Errorf("parsing %q: %w", path, err)
		}
		outpath := strings.TrimSuffix(path, extension) + ".go"
		bs, err := codegen.Generate(&codegen.Root{
			Package: filepath.Base(filepath.Dir(path)),
			Module:  *module,
		})
		if err != nil {
			return fmt.Errorf("generating from %q: %w", path, err)
		}
		if err := ioutil.WriteFile(outpath, bs, 0600); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		cmdutil.Fatalf("%v", err)
	}
}
