package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	types "github.com/joey00072/iza/pkg/types"
)

func RunContainer(args []string) error {

	container := args[0][:12]
	containerPath := filepath.Join(args[1], args[0])

	print("Container Path: ", containerPath, "\n")
	fmt.Println(">", args[2])

	syscall.Sethostname([]byte(container))
	syscall.Chroot(containerPath)
	syscall.Chdir("/")

	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("proc", 0)

	cmd := exec.Command(args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// image, err := ReadImage(filepath.Join(path, manifest.Config))
	// if err != nil {
	// 	return nil, fmt.Errorf("error reading image: %w", err)
	// }

	cmd.Run()

	return nil
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
