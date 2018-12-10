package helpers

import (
	"log"
	"os/exec"
	"syscall"
)

// RunCommandWithOutput runs the OS command and return its output. If the output
// returns error it also encapsulates it as a golang.error which is a return code
// of the command except zero
func RunCommandWithOutput(dir string, command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.Output()
	return string(output), err
}

// GetCommandStatus returns if we supposed to get return value as an int of a command
// this method can be used. It is practical when you use a command and process a
// failover acoording to a soecific return code
func GetCommandStatus(dir string, command string, args []string) (int, error) {
	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var err error
	// this time the execution is a little different
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
