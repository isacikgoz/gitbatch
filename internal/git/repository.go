package git

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
)

// Reference is the interface for commits, remotes and branches
type Reference interface {
	Next() *Reference
	Previous() *Reference
}

// Repository is the main entity of the application. The repository name is
// actually the name of its folder in the host's filesystem. It holds the go-git
// repository entity along with critic entities such as remote/branches and commits
type Repository struct {
	RepoID   string
	Name     string
	AbsPath  string
	ModTime  time.Time
	Repo     git.Repository
	Branches []*Branch
	Remotes  []*Remote
	Stasheds []*StashedItem
	State    *RepositoryState

	mutex     *sync.RWMutex
	listeners map[string][]RepositoryListener
}

// RepositoryState is the current pointers of a repository
type RepositoryState struct {
	workStatus WorkStatus
	Branch     *Branch
	Remote     *Remote
	Message    string
}

// RepositoryListener is a type for listeners
type RepositoryListener func(event *RepositoryEvent) error

// RepositoryEvent is used to transfer event-related data.
// It is passed to listeners when Publish() is called
type RepositoryEvent struct {
	Name string
	Data interface{}
}

// WorkStatus is the state of the repository for an operation
type WorkStatus struct {
	Status uint8
	Ready  bool
}

var (
	// Available implies repo is ready for the operation
	Available = WorkStatus{Status: 0, Ready: true}
	// Queued means repo is queued for a operation
	Queued = WorkStatus{Status: 1, Ready: false}
	// Working means an operation is just started for this repository
	Working = WorkStatus{Status: 2, Ready: false}
	// Paused is expected when a user interaction is required
	Paused = WorkStatus{Status: 3, Ready: true}
	// Success is the expected outcome of the operation
	Success = WorkStatus{Status: 4, Ready: true}
	// Fail is the unexpected outcome of the operation
	Fail = WorkStatus{Status: 5, Ready: false}
)

const (
	// RepositoryUpdated defines the topic for an updated repository.
	RepositoryUpdated = "repository.updated"
	// BranchUpdated defines the topic for an updated branch.
	BranchUpdated = "branch.updated"
)

// FastInitializeRepo initializes a Repository struct without its belongings.
func FastInitializeRepo(dir string) (r *Repository, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// get status of the file
	fstat, _ := f.Stat()
	rp, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	// initialize Repository with minimum viable fields
	r = &Repository{RepoID: RandomString(8),
		Name:    fstat.Name(),
		AbsPath: dir,
		ModTime: fstat.ModTime(),
		Repo:    *rp,
		State: &RepositoryState{
			workStatus: Available,
			Message:    "",
		},
		mutex:     &sync.RWMutex{},
		listeners: make(map[string][]RepositoryListener),
	}
	return r, nil
}

// InitializeRepo initializes a Repository struct with its belongings.
func InitializeRepo(dir string) (r *Repository, err error) {
	r, err = FastInitializeRepo(dir)
	if err != nil {
		return nil, err
	}
	// need nothing extra but loading additional components
	return r, r.loadComponents(true)
}

// loadComponents initializes the fields of a repository such as branches,
// remotes, commits etc. If reset, reload commit, remote pointers too
func (r *Repository) loadComponents(reset bool) error {
	if err := r.initRemotes(); err != nil {
		return err
	}

	if err := r.initBranches(); err != nil {
		return err
	}

	if err := r.SyncRemoteAndBranch(r.State.Branch); err != nil {
		return err
	}

	return r.loadStashedItems()
}

// Refresh the belongings of a repository, this function is called right after
// fetch/pull/merge operations
func (r *Repository) Refresh() error {
	var err error
	// error can be ignored since the file already exists when app is loading
	// if the Repository is only fast initialized, no need to refresh because
	// it won't contain its belongings
	if r.State.Branch == nil {
		return nil
	}
	file, _ := os.Open(r.AbsPath)
	fstat, _ := file.Stat()
	// re-initialize the go-git repository struct after supposed update
	rp, err := git.PlainOpen(r.AbsPath)
	if err != nil {
		return err
	}
	r.Repo = *rp
	// modification date may be changed
	r.ModTime = fstat.ModTime()
	if err := r.loadComponents(false); err != nil {
		return err
	}
	// we could send an event data but we don't need for this topic
	return r.Publish(RepositoryUpdated, nil)
}

// On adds new listener.
// listener is a callback function that will be called when event emits
func (r *Repository) On(event string, listener RepositoryListener) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// add listener to the specific event topic
	r.listeners[event] = append(r.listeners[event], listener)
}

// Publish publishes the data to a certain event by its name.
func (r *Repository) Publish(eventName string, data interface{}) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	// let's find listeners for this event topic
	listeners, ok := r.listeners[eventName]
	if !ok {
		return nil
	}
	// now notify the listeners and channel the data
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

// WorkStatus returns the state of the repository such as queued, failed etc.
func (r *Repository) WorkStatus() WorkStatus {
	return r.State.workStatus
}

// SetWorkStatus sets the state of repository and sends repository updated event
func (r *Repository) SetWorkStatus(ws WorkStatus) {
	r.State.workStatus = ws
	// we could send an event data but we don't need for this topic
	_ = r.Publish(RepositoryUpdated, nil)
}

func (r *Repository) String() string {
	return r.Name
}

func Create(dir string) (*Repository, error) {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return InitializeRepo(dir)
}
