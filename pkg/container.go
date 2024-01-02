package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	types "github.com/joey00072/iza/pkg/types"
)

func RunContainer(args []string) error {

	container := args[0][:12]
	containerPath := filepath.Join(args[1], args[0])

	// print("Container Path: ", containerPath, "\n")

	imgagePath, err := addPrefixToLastDir(containerPath, "sha256:")
	if err != nil {
		return fmt.Errorf("error adding prefix to last dir: %w", err)
	}

	image, err := ReadImage(imgagePath)
	if err != nil {
		return fmt.Errorf("error reading image: %w", err)
	}

	syscall.Sethostname([]byte(container))
	syscall.Chroot(containerPath)
	syscall.Chdir("/")

	for _, env := range image.Config.Env {
		kv := strings.Split(env, "=")
		syscall.Setenv(kv[0], kv[1])
	}

	time.Sleep(time.Second * 4)

	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("proc", 0)

	fmt.Println("Running Container...")
	fmt.Println("---------------------")

	var cmd *exec.Cmd
	if len(args) > 2 {
		cmd = exec.Command(args[2])
	} else {
		var cammands []string
		if image.ContainerConfig != nil {
			cammands = image.ContainerConfig.Cmd

		} else if image.Config != nil {
			cammands = image.Config.Cmd
		}
		cmd = exec.Command(cammands[0], cammands[1:]...)
	}

	// cmd = exec.Command(args[2])

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return nil
}

func addPrefixToLastDir(path, prefix string) (string, error) {
	parts := strings.Split(filepath.ToSlash(path), "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("path does not have enough components")
	}

	parts[len(parts)-1] = prefix + parts[len(parts)-1]

	return strings.Join(parts, "/"), nil
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
