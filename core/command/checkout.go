package command

import "github.com/isacikgoz/gitbatch/core/git"

// CheckoutOptions defines the rules of checkout command
type CheckoutOptions struct {
	TargetRef      string
	CreateIfAbsent bool
	CommandMode    Mode
}

// Checkout is a wrapper function for "git checkout" command.
func Checkout(r *git.Repository, o CheckoutOptions) error {
	return nil
}
