package command

import (
	"errors"
	"os/exec"
	"strings"
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

// Return runs the OS command and returns its exit code. It is practical when
// you use a command and process a failover according to a specific return
// code. If the command could not be run, the exit code is -1.
func Return(d string, c string, args []string) (int, error) {
	cmd := exec.Command(c, args...)
	if d != "" {
		cmd.Dir = d
	}
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// the program has exited with an exit code != 0
		return exitErr.ExitCode(), err
	}
	if err != nil {
		return -1, err
	}
	return 0, nil
}

// trimTrailingNewline removes the trailing new line form a string. this method
// is used mostly on outputs of a command
func trimTrailingNewline(s string) string {
	if strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\r") {
		return s[:len(s)-1]
	}
	return s
}
