package git

import (
	"log"
	"os/exec"
	"strings"
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
	output, err := cmd.CombinedOutput()
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

// TrimTrailingNewline removes the trailing new line form a string. this method
// is used mostly on outputs of a command
func TrimTrailingNewline(str string) string {
	if strings.HasSuffix(str, "\n") {
		return str[:len(str)-1]
	}
	return str
}

// GenericGitCommand runs any git command without expecting output
func GenericGitCommand(repoPath string, args []string) error {
	_, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return err
	}
	return nil
}

// GenericGitCommandWithOutput runs any git command with returning output
func GenericGitCommandWithOutput(repoPath string, args []string) (string, error) {
	out, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return out, err
	}
	return TrimTrailingNewline(out), nil
}

// GenericGitCommandWithErrorOutput runs any git command with returning output
func GenericGitCommandWithErrorOutput(repoPath string, args []string) (string, error) {
	out, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return TrimTrailingNewline(out), err
	}
	return TrimTrailingNewline(out), nil
}

// GitShow is conventional git show command without any argument
func GitShow(repoPath, hash string) string {
	args := []string{"show", hash}
	diff, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// GitShowEmail gets author's e-mail with git show command
func GitShowEmail(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ae", hash}
	diff, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// GitShowBody gets body of the commit with git show
func GitShowBody(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%B", hash}
	diff, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return err.Error()
	}
	return diff
}

// GitShowDate gets commit's date with git show as string
func GitShowDate(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ai", hash}
	diff, err := RunCommandWithOutput(repoPath, "git", args)
	if err != nil {
		return "?"
	}
	return diff
}

// StatusWithGit returns the plaintext short status of the repo
func (e *RepoEntity) StatusWithGit() string {
	args := []string{"status"}
	status, err := RunCommandWithOutput(e.AbsPath, "git", args)
	if err != nil {
		return "?"
	}
	return status
}
