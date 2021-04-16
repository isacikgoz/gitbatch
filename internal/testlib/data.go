package testlib

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:embed test-data-master.zip
var zipFile embed.FS

func ExtractTestRepository(dir string) (string, error) {
	data, err := zipFile.ReadFile("test-data-master.zip")
	if err != nil {
		return "", fmt.Errorf("could not extract test data: %w", err)
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("could not read test data: %w", err)
	}

	for _, f := range r.File {
		if err := extractFile(dir, f); err != nil {
			return "", fmt.Errorf("could not extract test data: %w", err)
		}
	}

	return filepath.Join(dir, "test-data"), nil
}

func extractFile(dst string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(dst, f.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dst)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(path, f.Mode())
	}

	err = os.MkdirAll(filepath.Dir(path), f.Mode())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, rc)
	if err != nil {
		return err
	}

	return nil
}
