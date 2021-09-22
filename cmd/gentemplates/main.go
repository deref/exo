package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/deref/exo/internal/template"
	"github.com/go-git/go-git/v5"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := genTemplates(); err != nil {
		panic(err)
	}
}

func genTemplates() error {
	dir, err := os.MkdirTemp("", "exo-template-clone-")
	if err != nil {
		return fmt.Errorf("making temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	outDir, err := os.MkdirTemp("", "exo-templates-")
	if err != nil {
		return fmt.Errorf("making temp dir: %w", err)
	}
	fmt.Println(outDir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:          "https://github.com/railwayapp/starters",
		SingleBranch: true,
		Tags:         git.NoTags,
		Depth:        1,
	})
	if err != nil {
		return fmt.Errorf("cloning repo: %w", err)
	}

	examplesPath := path.Join(dir, "examples")
	entries, err := ioutil.ReadDir(examplesPath)
	if err != nil {
		return fmt.Errorf("reading examples dir: %w", err)
	}

	ctx := context.Background()
	var eg errgroup.Group
	for _, e := range entries {
		if e.Name() == "node_modules" {
			continue
		}

		if e.IsDir() {
			templateName := e.Name()
			templateDir := path.Join(examplesPath, templateName)

			tmplOutDir := path.Join(outDir, templateName)
			if err := os.Mkdir(tmplOutDir, 0750); err != nil {
				return fmt.Errorf("creating template dir: %w", err)
			}

			eg.Go(func() error {
				if err := template.MakeTemplateFiles(ctx, templateDir, tmplOutDir); err != nil {
					return fmt.Errorf("making template files: %w", err)
				}
				return nil
			})
		}
	}

	return eg.Wait()
}
