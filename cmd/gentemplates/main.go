package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := genTemplates(); err != nil {
		panic(err)
	}
}

// From https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
func compress(src string, buf io.Writer) error {
	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}

func genTemplates() error {
	dir, err := os.MkdirTemp("", "exo-template-clone-")
	if err != nil {
		return fmt.Errorf("making temp dir: %w", err)
	}

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

	var eg errgroup.Group
	for _, e := range entries {
		if e.Name() == "node_modules" {
			continue
		}

		if e.IsDir() {
			templateName := e.Name()
			templateDir := path.Join(examplesPath, templateName)
			tarFile := templateDir

			if err := os.Mkdir(path.Join(outDir, templateName), 0750); err != nil {
				return fmt.Errorf("creating template dir: %w", err)
			}

			f, err := os.Create(path.Join(outDir, templateName, "files.tar.gz"))
			if err != nil {
				return fmt.Errorf("failed to open file %q, %v", tarFile, err)
			}
			defer f.Close()

			eg.Go(func() error {
				return compress(templateDir, f)
			})
		}
	}

	return eg.Wait()
}
