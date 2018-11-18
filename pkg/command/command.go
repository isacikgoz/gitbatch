package command

import (
	"os/exec"
	"syscall"
	"log"
)

func RunCommandWithOutput(dir string, command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.Output() 
	return string(output), err
}

func GetCommandStatus(dir string, command string, args []string) (int, error) {
	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var err error
	if err := cmd.Start(); err != nil {
         return -1, err
    }
    if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0

            // This works on both Unix and Windows. Although package
            // syscall is generally platform dependent, WaitStatus is
            // defined for both Unix and Windows and in both cases has
            // an ExitStatus() method with the same signature.
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                statusCode := status.ExitStatus()
        		return statusCode, err
            }
        } else {
            log.Fatalf("cmd.Wait: %v", err)
        }
    }
    return -1, err
}
