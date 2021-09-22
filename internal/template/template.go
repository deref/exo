package template

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var archiveFolderName = "files"
var tarName = archiveFolderName + ".tar.gz"

// GetTemplateFiles returns the path to a directory that contains the intial
// project template.
func GetTemplateFiles(ctx context.Context, templateURL string) (string, error) {
	dir, err := os.MkdirTemp("", "exo-template-clone-")
	if err != nil {
		return "", fmt.Errorf("making temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	url := fmt.Sprintf("%s/%s", templateURL, tarName)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("getting template files: %w", err)
	}
	defer resp.Body.Close()

	err = uncompress(resp.Body, dir)
	if err != nil {
		return "", fmt.Errorf("untarring template: %w", err)
	}
	return filepath.Join(dir, archiveFolderName), nil
}

// MakeTemplateFiles takes a directory of files to template and outputs the
// resultant template into the given directory.
func MakeTemplateFiles(ctx context.Context, inputDir, outputDir string) error {
	tarFile := filepath.Join(outputDir, tarName)
	f, err := os.Create(tarFile)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", tarFile, err)
	}
	defer f.Close()

	if err := compress(inputDir, f); err != nil {
		return fmt.Errorf("compressing template files: %w", err)
	}
	return nil
}
