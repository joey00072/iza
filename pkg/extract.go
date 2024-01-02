package pkg

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	types "github.com/joey00072/iza/pkg/types"
	utils "github.com/joey00072/iza/pkg/utils"
)

func ExtractImage(imageOpt types.ImageOptions) (*types.Container, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %w", err)
	}
	path := imageOpt.CachePath
	if path == "" {
		path = filepath.Join(cwd, types.DefaultCachePath)
	}

	path = filepath.Join(path, imageOpt.ImageName)

	createIfdontExist := true
	err = utils.EnsurePathExists(path, createIfdontExist)
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
		tarGzFile := filepath.Join(path, layer)
		exec.Command("tar", "-C", containerDir, "-xzf", tarGzFile).Run()
		// if err != nil {
		// 	fmt.Println("Error during extraction:", err)
		// 	return nil, fmt.Errorf("error extracting tarball: %w", err)
		// }

	}

	conatiner := &types.Container{
		ID:  conatinerName,
		Dir: path,
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
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
