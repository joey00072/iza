package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	types "github.com/joey00072/iza/pkg/types"
	"github.com/joey00072/iza/pkg/utils"
)

func Run(imageOpt types.ImageOptions, args []string) error {

	if !utils.TarExits(imageOpt.TarballPath) {
		err := PullImage(imageOpt)
		if err != nil {
			return fmt.Errorf("error pulling image: %w", err)
		}
	}

	conatiner, err := ExtractImage(imageOpt)
	if err != nil {
		return fmt.Errorf("error extracting image: %w", err)
	}
	defer utils.DeleteDir(conatiner.Dir)

	containerArgs := []string{"cprocess"}
	containerArgs = append(containerArgs, conatiner.ID)
	containerArgs = append(containerArgs, conatiner.Dir)
	containerArgs = append(containerArgs, args[1:]...)

	cmd := exec.Command("/proc/self/exe", containerArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	cmd.Run()

	return nil
}
