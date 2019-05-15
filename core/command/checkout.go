package command

import "github.com/isacikgoz/gitbatch/core/git"

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
	if branch == nil && o.CreateIfAbsent {

	}
	var msg string
	if branch != nil {
		r.SetWorkStatus(git.Success)
		msg = "switched to " + o.TargetRef
		if err := r.Checkout(branch); err != nil {
			msg = err.Error()
		}
	}
	r.State.Message = msg
	return r.Refresh()
}
