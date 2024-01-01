package pkg

import (
	"fmt"

	types "github.com/joey00072/iza/pkg/types"
	"github.com/joey00072/iza/pkg/utils"
)

func RunContainer(imageOpt types.ImageOptions, args []string) error {
	fmt.Printf("Tarpath %v \n", imageOpt.TarballPath)
	if !utils.TarExits(imageOpt.TarballPath) {
		PullImage(imageOpt)
	}
	fmt.Println("ckpt 1")

	conatiner, err := ExtractImage(imageOpt)
	if err != nil {
		return fmt.Errorf("error extracting image: %w", err)
	}
	defer utils.DeleteDir(conatiner.ImageDir)

	fmt.Println("Done Running Container")
	return nil
}
