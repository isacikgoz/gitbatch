package command

import (
	"os/exec"

	"github.com/isacikgoz/gitbatch/internal/git"
)

// CheckoutOptions defines the rules of checkout command
type CheckoutOptions struct {
	TargetRef      string
	CreateIfAbsent bool
	CommandMode    Mode
}

// Checkout is a wrapper function for "git checkout" command.
func Checkout(r *git.Repository, o *CheckoutOptions) error {
	var branch *git.Branch
	for _, b := range r.Branches {
		if b.Name == o.TargetRef {
			branch = b
			break
		}
	}
	msg := "checkout in progress"
	if branch != nil {
		if err := r.Checkout(branch); err != nil {
			r.SetWorkStatus(git.Fail)
			msg = err.Error()
		} else {
			r.SetWorkStatus(git.Success)
			msg = "switched to " + o.TargetRef
		}
	} else if o.CreateIfAbsent {
		args := []string{"checkout", "-b", o.TargetRef}
		cmd := exec.Command("git", args...)
		cmd.Dir = r.AbsPath
		_, err := cmd.CombinedOutput()
		if err != nil {
			r.SetWorkStatus(git.Fail)
			msg = err.Error()
		} else {
			r.SetWorkStatus(git.Success)
			msg = "switched to " + o.TargetRef
		}
	}
	r.State.Message = msg
	return r.Refresh()
}
