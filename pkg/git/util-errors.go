package git

import (
	"errors"
)

var (
	// ErrGitCommand is thrown when git command returned an error code
	ErrGitCommand = errors.New("Git command returned error code")
	// ErrAuthenticationRequired is thrown when an authentication required on
	// a remote operation
	ErrAuthenticationRequired = errors.New("Authentication required")
	// ErrAuthorizationFailed is thrown when authorization failed while trying
	// to authenticate with remote
	ErrAuthorizationFailed = errors.New("Authorization failed")
	// ErrInvalidAuthMethod is thrown when invalid auth method is invoked
	ErrInvalidAuthMethod = errors.New("invalid auth method")
	// ErrAlreadyUpToDate is thrown when a repository is already up to date
	// with its src on merge/fetch/pull
	ErrAlreadyUpToDate = errors.New("Already up to date")
	// ErrCouldNotFindRemoteRef is thrown when trying to fetch/pull cannot
	// find suitable remote reference
	ErrCouldNotFindRemoteRef = errors.New("Could not find remote ref")
)
