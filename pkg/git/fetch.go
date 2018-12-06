package git

import (
	log "github.com/sirupsen/logrus"
)

var fetchCommand = "fetch"

type FetchOptions struct {
    // Name of the remote to fetch from. Defaults to origin.
    RemoteName string
    // Before fetching, remove any remote-tracking references that no longer
    // exist on the remote.
    Prune bool
    // Show what would be done, without making any changes.
    DryRun bool
    // Force allows the fetch to update a local branch even when the remote
    // branch does not descend from it.
    Force bool
}

// Fetch branches refs from one or more other repositories, along with the
// objects necessary to complete their histories
func Fetch(entity *RepoEntity, options FetchOptions) error {
	args := make([]string, 0)
	args = append(args, fetchCommand)
	if len(options.RemoteName) > 0 {
		args = append(args, options.RemoteName)
	}
	if options.Prune {
		args = append(args, "-p")
	}
	if options.Force {
		args = append(args, "-f")
	}
	if options.DryRun {
		args = append(args, "--dry-run")
	}
	if err := GenericGitCommand(entity.AbsPath, args); err != nil {
		log.Warn("Error while fetching")
		return err
	}
	return nil
}