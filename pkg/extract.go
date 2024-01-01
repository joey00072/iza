package pkg

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	types "github.com/joey00072/iza/pkg/types"
	utils "github.com/joey00072/iza/pkg/utils"
)

func ExtractImage(imageOpt types.ImageOptions) error {
	path := imageOpt.CachePath
	if path == "" {
		path = types.DefaultCachePath

		return nil
	}
	createIfdontExist := true
	err := utils.EnsurePathExists(path, createIfdontExist)
	if err != nil {
		return fmt.Errorf("error ensuring path exists: %w", err)
	}

	tar := imageOpt.TarballPath

	err = ExtractTar(tar, path)
	if err != nil {
		return fmt.Errorf("error extracting tarball: %w", err)
	}

	return nil
}

// ExtractTar extracts a tar file to a specified destination.
func ExtractTar(tarPath, destPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar file: %w", err)
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to make directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		}
	}

	return nil
}
