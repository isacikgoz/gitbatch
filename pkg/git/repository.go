package git

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/isacikgoz/gitbatch/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

// RepoEntity is the main entity of the application. The repository name is
// actually the name of its folder in the host's filesystem. It holds the go-git
// repository entity along with critic entites such as remote/branches and commits
type RepoEntity struct {
	RepoID     string
	Name       string
	AbsPath    string
	ModTime    time.Time
	Repository git.Repository
	Branch     *Branch
	Branches   []*Branch
	Remote     *Remote
	Remotes    []*Remote
	Commit     *Commit
	Commits    []*Commit
	Stasheds   []*StashedItem
	state      RepoState

	mutex     *sync.RWMutex
	listeners map[string][]RepositoryListener
}

// RepositoryListener is a type for listeners
type RepositoryListener func(event *RepositoryEvent) error

// RepositoryEvent is used to transfer event-related data.
// It is passed to listeners when Emit() is called
type RepositoryEvent struct {
	Name string
	Data interface{}
}

// RepoState is the state of the repository for an operation
type RepoState uint8

const (
	// Available implies repo is ready for the operation
	Available RepoState = 0
	// Queued means repo is queued for a operation
	Queued RepoState = 1
	// Working means an operation is just started for this repository
	Working RepoState = 2
	// Paused is expected when a user interaction is required
	Paused RepoState = 3
	// Success is the expected outcome of the operation
	Success RepoState = 4
	// Fail is the unexpected outcome of the operation
	Fail RepoState = 5

	// This is the repository updated topic
	RepositoryUpdated = "repository.updated"
)

// InitializeRepo initializes a RepoEntity struct with its belongings.
func InitializeRepo(directory string) (entity *RepoEntity, err error) {
	entity, err = FastInitializeRepo(directory)
	if err != nil {
		return nil, err
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
	if err = entity.loadStashedItems(); err != nil {
		// TODO: fix here.
	}
	if len(entity.Remotes) > 0 {
		// TODO: tend to take origin/master as default
		entity.Remote = entity.Remotes[0]
		if entity.Branch == nil {
			return nil, errors.New("Unable to find a valid branch")
		}
		if err = entity.Remote.SyncBranches(entity.Branch.Name); err != nil {
			// probably couldn't find, but its ok.
		}
	} else {
		// if there is no remote, this project is totally useless actually
		return entity, errors.New("There is no remote for this repository: " + directory)
	}
	return entity, nil
}

// FastInitializeRepo initializes a RepoEntity struct without its belongings.
func FastInitializeRepo(directory string) (entity *RepoEntity, err error) {
	file, err := os.Open(directory)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": directory,
		}).Trace("Cannot open as directory")
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(directory)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": directory,
		}).Trace("Cannot open directory as a git repository")
		return nil, err
	}
	// initialize entity with minimum viable fields
	entity = &RepoEntity{RepoID: helpers.RandomString(8),
		Name:       fileInfo.Name(),
		AbsPath:    directory,
		ModTime:    fileInfo.ModTime(),
		Repository: *r,
		state:      Available,
		mutex:      &sync.RWMutex{},
		listeners:  make(map[string][]RepositoryListener),
	}
	return entity, nil
}

// On adds new listener.
// listener is a callback function that will be called when event emits
func (entity *RepoEntity) On(event string, listener RepositoryListener) {
	entity.mutex.Lock()
	defer entity.mutex.Unlock()

	entity.listeners[event] = append(entity.listeners[event], listener)
}

// Emit notifies listeners about the event
func (entity *RepoEntity) Emit(eventName string, data interface{}) error {
	entity.mutex.RLock()
	defer entity.mutex.RUnlock()

	listeners, ok := entity.listeners[eventName]
	if !ok {
		return nil
	}

	for i := range listeners {
		event := &RepositoryEvent{
			Name: eventName,
			Data: data,
		}
		if err := listeners[i](event); err != nil {
			return err
		}
	}
	return nil
}

// State returns the state of the repository such as queued, failed etc.
func (entity *RepoEntity) State() RepoState {
	return entity.state
}

// SetState sets the state of repository and sends repository updated event
func (entity *RepoEntity) SetState(state RepoState) {
	entity.state = state
	// we could send an event data but we don't need for this topic
	entity.Emit(RepositoryUpdated, nil)
}

// Refresh the belongings of a repositoriy, this function is called right after
// fetch/pull/merge operations
func (entity *RepoEntity) Refresh() error {
	var err error
	// error can be ignored since the file already exists when app is loading
	if entity.Branch == nil {
		return nil
	}
	file, _ := os.Open(entity.AbsPath)
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(entity.AbsPath)
	if err != nil {
		return err
	}
	entity.Repository = *r
	entity.ModTime = fileInfo.ModTime()
	if err := entity.loadLocalBranches(); err != nil {
		return err
	}
	entity.Branch.Clean = entity.isClean()
	entity.RefreshPushPull()
	if err := entity.loadCommits(); err != nil {
		return err
	}
	if err := entity.loadRemotes(); err != nil {
		return err
	}
	err = entity.loadStashedItems()
	// we could send an event data but we don't need for this topic
	return entity.Emit(RepositoryUpdated, nil)
}
