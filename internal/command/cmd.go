package command

import (
	"log"
	"os/exec"
	"strings"
	"syscall"
)

// Mode indicates that whether command should run native code or use git
// command to operate.
type Mode uint8

const (
	// ModeLegacy uses traditional git command line tool to operate
	ModeLegacy = iota
	// ModeNative uses native implementation of given git command
	ModeNative
)

// Run runs the OS command and return its output. If the output
// returns error it also encapsulates it as a golang.error which is a return code
// of the command except zero
func Run(d string, c string, args []string) (string, error) {
	cmd := exec.Command(c, args...)
	if d != "" {
		cmd.Dir = d
	}
	output, err := cmd.CombinedOutput()
	return trimTrailingNewline(string(output)), err
}

// Return returns if we supposed to get return value as an int of a command
// this method can be used. It is practical when you use a command and process a
// failover according to a specific return code
func Return(d string, c string, args []string) (int, error) {
	cmd := exec.Command(c, args...)
	if d != "" {
		cmd.Dir = d
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

// trimTrailingNewline removes the trailing new line form a string. this method
// is used mostly on outputs of a command
func trimTrailingNewline(s string) string {
	if strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\r") {
		return s[:len(s)-1]
	}
	return s
}
