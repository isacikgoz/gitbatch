package git

import (
	"errors"
	"os"

	"github.com/isacikgoz/gitbatch/pkg/helpers"
	"gopkg.in/src-d/go-git.v4"
)

// the main entity of the application. The repository name is actually the name
// of its folder in the host's filesystem. It holds the go-git repository entity
// along with critic entites such as remote/branches and commits
type RepoEntity struct {
	RepoID     string
	Name       string
	AbsPath    string
	Repository git.Repository
	Branch     *Branch
	Branches   []*Branch
	Remote     *Remote
	Remotes    []*Remote
	Commit     *Commit
	Commits    []*Commit
	State      RepoState
}

// it is the state of the repository for an operation
type RepoState uint8

const (
	Available RepoState = 0
	Queued    RepoState = 1
	Working   RepoState = 2
	Success   RepoState = 3
	Fail      RepoState = 4
)

// initializee a RepoEntity struct with its belongings.
func InitializeRepository(directory string) (entity *RepoEntity, err error) {
	file, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}
	entity = &RepoEntity{RepoID: helpers.RandomString(8),
		Name:       fileInfo.Name(),
		AbsPath:    directory,
		Repository: *r,
		State:      Available,
	}
	// after we intiate the struct we can fill its values
	entity.loadLocalBranches()
	entity.loadCommits()
	// handle if there is no commit, maybe?
	if len(entity.Commits) > 0 {
		// select first commit
		entity.Commit = entity.Commits[0]
	} else {
		return entity, errors.New("There is no commit for this repository: " + directory)
	}
	// lets load remotes this time
	entity.loadRemotes()
	// set the active branch to repositories HEAD
	entity.Branch = entity.getActiveBranch()
	if len(entity.Remotes) > 0 {
		// TODO: tend to take origin/master as default
		entity.Remote = entity.Remotes[0]
		// TODO: same code on 3 different occasion, maybe something wrong?
		if err = entity.Remote.switchRemoteBranch(entity.Remote.Name + "/" + entity.Branch.Name); err != nil {
			// probably couldn't find, but its ok.
		}
	} else {
		// if there is no remote, this project is totally useless actually
		return entity, errors.New("There is no remote for this repository: " + directory)
	}
	return entity, nil
}

// Incorporates changes from a remote repository into the current branch. In
// its default mode, git pull is shorthand for git fetch followed by git merge
// <branch>
func (entity *RepoEntity) Pull() error {
	// TODO: Migrate this code to src-d/go-git
	// 2018-11-25: tried but it fails, will investigate.
	rm := entity.Remote.Name
	if err := entity.FetchWithGit(rm); err != nil {
		return err
	}
	entity.Checkout(entity.Branch)
	if err := entity.MergeWithGit(entity.Remote.Branch.Name); err != nil {
		entity.Refresh()
		return err
	}
	entity.Refresh()
	entity.Checkout(entity.Branch)
	return nil
}

// Fetch branches refs from one or more other repositories, along with the
// objects necessary to complete their histories
func (entity *RepoEntity) Fetch() error {
	rm := entity.Remote.Name
	if err := entity.FetchWithGit(rm); err != nil {
		return err
	}
	entity.Refresh()
	entity.Checkout(entity.Branch)
	return nil
}

// Incorporates changes from the named commits or branches into the current
// branch
func (entity *RepoEntity) Merge() error {
	entity.Checkout(entity.Branch)
	if err := entity.MergeWithGit(entity.Remote.Branch.Name); err != nil {
		entity.Refresh()
		return err
	}
	entity.Refresh()
	return nil
}

// refresh the belongings of a repositoriy, this function is called right after
// fetch/pull/merge operations
func (entity *RepoEntity) Refresh() error {
	r, err := git.PlainOpen(entity.AbsPath)
	if err != nil {
		return err
	}
	entity.Repository = *r
	if err := entity.loadLocalBranches(); err != nil {
		return err
	}
	if err := entity.loadCommits(); err != nil {
		return err
	}
	if err := entity.loadRemotes(); err != nil {
		return err
	}
	return nil
}
