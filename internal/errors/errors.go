package errors

import (
	"strings"
)

// GitError is the errors from git package
type GitError string

const (
	// ErrGitCommand is thrown when git command returned an error code
	ErrGitCommand GitError = "git command returned error code"
	// ErrAuthenticationRequired is thrown when an authentication required on
	// a remote operation
	ErrAuthenticationRequired GitError = "authentication required"
	// ErrAuthorizationFailed is thrown when authorization failed while trying
	// to authenticate with remote
	ErrAuthorizationFailed GitError = "authorization failed"
	// ErrInvalidAuthMethod is thrown when invalid auth method is invoked
	ErrInvalidAuthMethod GitError = "invalid auth method"
	// ErrAlreadyUpToDate is thrown when a repository is already up to date
	// with its src on merge/fetch/pull
	ErrAlreadyUpToDate GitError = "already up to date"
	// ErrCouldNotFindRemoteRef is thrown when trying to fetch/pull cannot
	// find suitable remote reference
	ErrCouldNotFindRemoteRef GitError = "could not find remote ref"
	// ErrMergeAbortedTryCommit indicates that the repositort is not clean and
	// some changes may conflict with the merge
	ErrMergeAbortedTryCommit GitError = "stash/commit changes. aborted"
	// ErrRemoteBranchNotSpecified means that default remote branch is not set
	// for the current branch. can be setted with "git config --local --add
	// branch.<your branch name>.remote=<your remote name> "
	ErrRemoteBranchNotSpecified GitError = ("upstream not set")
	// ErrRemoteNotFound is thrown when the remote is not reachable. It may be
	// caused by the deletion of the remote or connectivity problems
	ErrRemoteNotFound GitError = ("remote not found")
	// ErrConflictAfterMerge is thrown when a conflict occurs at merging two
	// references
	ErrConflictAfterMerge GitError = ("conflict while merging")
	// ErrUnmergedFiles possibly occurs after a conflict
	ErrUnmergedFiles GitError = ("unmerged files detected")
	// ErrReferenceBroken thrown when unable to resolve reference
	ErrReferenceBroken GitError = ("unable to resolve reference")
	// ErrPermissionDenied is thrown when ssh authentication occurs
	ErrPermissionDenied GitError = ("permission denied")
	// ErrOverwrittenByMerge is the thrown when there is un-tracked files on working tree
	ErrOverwrittenByMerge GitError = ("move or remove un-tracked files before merge")
	// ErrUserEmailNotSet is thrown if there is no configured user email while
	// commit command
	ErrUserEmailNotSet GitError = ("user email not configured")
	// ErrUnclassified is unconsidered error type
	ErrUnclassified GitError = ("unclassified error")
	// NoErrIterationHalted is thrown for catching stops in interators
	NoErrIterationHalted GitError = ("iteration halted")
)

func (e GitError) Error() string {
	return string(e)
}

// ParseGitError takes git output as an input and tries to find some meaningful
// errors can be used by the app
func ParseGitError(out string, err error) error {
	if strings.Contains(out, "error: Your local changes to the following files would be overwritten by merge") {
		return ErrMergeAbortedTryCommit
	} else if strings.Contains(out, "ERROR: Repository not found") {
		return ErrRemoteNotFound
	} else if strings.Contains(out, "for your current branch, you must specify a branch on the command line") {
		return ErrRemoteBranchNotSpecified
	} else if strings.Contains(out, "Automatic merge failed; fix conflicts and then commit the result") {
		return ErrConflictAfterMerge
	} else if strings.Contains(out, "error: Pulling is not possible because you have unmerged files.") {
		return ErrUnmergedFiles
	} else if strings.Contains(out, "unable to resolve reference") {
		return ErrReferenceBroken
	} else if strings.Contains(out, "git config --global add user.email") {
		return ErrUserEmailNotSet
	} else if strings.Contains(out, "Permission denied (publickey)") {
		return ErrPermissionDenied
	} else if strings.Contains(out, "would be overwritten by merge") {
		return ErrOverwrittenByMerge
	}
	return ErrUnclassified
}
