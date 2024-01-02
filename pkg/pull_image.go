package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/cache"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	types "github.com/joey00072/iza/pkg/types"
)

// PullImage pulls a docker image from a registry and saves it as a tarball.
// It reports progress and returns an error if the image cannot be pulled or saved.
func PullImage(imageOpt types.ImageOptions) error {
	imageName := imageOpt.ImageName
	tarballPath := imageOpt.TarballPath
	cachePath := imageOpt.CachePath // Use the provided CachePath

	imageMap := map[string]v1.Image{}

	// Configure crane options
	ctx := context.Background()
	o := crane.Options{
		Remote: []remote.Option{
			remote.WithContext(ctx),
		},
	}

	ref, err := name.ParseReference(imageName, o.Name...)
	if err != nil {
		return fmt.Errorf("error parsing image %s: %w", imageName, err)
	}

	fmt.Println("Connecting to remote registry...")

	rmt, err := remote.Get(ref, o.Remote...)
	if err != nil {
		return fmt.Errorf("error getting remote image %s: %w", imageName, err)
	}

	fmt.Println("Image found, starting download...")

	img, err := rmt.Image()
	if err != nil {
		return fmt.Errorf("error getting image from remote %s: %w", imageName, err)
	}

	fmt.Println("Download complete.")

	if cachePath != "" {
		img = cache.Image(img, cache.NewFilesystemCache(cachePath))
	}

	imageMap[imageName] = img

	tarballPath = strings.Replace(tarballPath, "/", "-", -1)

	// Pull the image from the registry and save it as a tarball
	if err := crane.MultiSave(imageMap, tarballPath); err != nil {
		return fmt.Errorf("error saving image %s as a tarball: %w", imageName, err)
	}

	fmt.Printf("Successfully saved %s as a tarball at %s\n", imageName, tarballPath)
	return nil
}

func createFileIfNotExists(filePath string) error {
	// Attempt to create the file. If the file already exists, an error will be returned.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
