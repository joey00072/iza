package pkg

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	types "github.com/joey00072/iza/pkg/types"
	utils "github.com/joey00072/iza/pkg/utils"
)

func ExtractImage(imageOpt types.ImageOptions) (*types.Container, error) {
	path := imageOpt.CachePath
	if path == "" {
		path = types.DefaultCachePath
	}

	path = filepath.Join(path, imageOpt.ImageName)

	createIfdontExist := true
	err := utils.EnsurePathExists(path, createIfdontExist)
	if err != nil {
		return nil, fmt.Errorf("error ensuring path exists: %w", err)
	}

	tar := imageOpt.TarballPath

	err = ExtractTar(tar, path, false)
	if err != nil {
		return nil, fmt.Errorf("error extracting tarball: %w", err)
	}

	manifestPath := filepath.Join(path, "manifest.json")
	manifest, err := ReadManifest(manifestPath)

	conatinerName := strings.Replace(manifest.Config, "sha256:", "", -1)
	fmt.Println("Conatiner Name:", conatinerName)
	containerDir := filepath.Join(path, conatinerName)
	containerExits, err := utils.DirExists(containerDir)
	if err != nil {
		return nil, fmt.Errorf("error checking container dir: %w", err)
	}

	if containerExits {
		fmt.Printf("Container alread exits at %v \n", conatinerName)
		return nil, fmt.Errorf("Container Dir alread exits at %v", containerDir)
	}

	err = utils.EnsurePathExists(containerDir, createIfdontExist)
	if err != nil {
		return nil, fmt.Errorf("Container Dir alread exits at %v : %w", containerDir, err)
	}

	for _, layer := range manifest.Layers {
		targz := filepath.Join(path, layer)
		fmt.Printf("targz: %v \n", targz)
		err = ExtractTar(targz, containerDir, true)
		if err != nil {
			return nil, fmt.Errorf("error extracting tarball: %w", err)
		}
	}

	image, err := ReadImage(filepath.Join(path, manifest.Config))
	if err != nil {
		return nil, fmt.Errorf("error reading image: %w", err)
	}

	conatiner := &types.Container{
		ID:       conatinerName,
		Dir:      containerDir,
		Image:    image,
		ImageDir: path,
	}

	return conatiner, nil
}

func ReadManifest(filename string) (*types.Manifest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var manifest []types.Manifest
	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest[0], nil
}

func ReadImage(filename string) (*types.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var image types.Image
	err = json.Unmarshal(bytes, &image)
	if err != nil {
		return nil, err
	}

	return &image, nil

}

func ExtractTar(tarPath, destPath string, gz bool) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer file.Close()

	var tarReader *tar.Reader
	if gz {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		tarReader = tar.NewReader(gzipReader)
	} else {
		tarReader = tar.NewReader(file)
	}

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
