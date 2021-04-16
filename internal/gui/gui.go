package gui

import (
	"fmt"
	"sort"
	"sync"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/job"
	"github.com/isacikgoz/gitbatch/internal/load"
	"github.com/jroimartin/gocui"
)

// Gui struct hold the gocui struct along with the gui's state, also keybindings
// are tied with this struct in order to render those in different occasions
type Gui struct {
	g           *gocui.Gui
	KeyBindings []*KeyBinding
	State       guiState
	mutex       *sync.Mutex
	order       Layout
}

// guiState struct holds the repositories, directories, mode and queue of the
// gui object. These values are not static
type guiState struct {
	Repositories  []*git.Repository
	Directories   []string
	Mode          mode
	Queue         *job.Queue
	FailoverQueue *job.Queue
	targetBranch  string
	totalBranches []*branchCountMap
}

// this struct encapsulates the name and title of a view. the name of a view is
// passed around so much it is added so that I don't need to write names again
type viewFeature struct {
	Name  string
	Title string
}

// mode of the gui
type mode struct {
	ModeID        ModeID
	DisplayString string
	CommandString string
}

type branchCountMap struct {
	BranchName string
	Count      int
}

// ModeID is the mode indicator for the gui
type ModeID string

// Layout tells the gui how to order views
type Layout int

const (
	// FetchMode puts the gui in fetch state
	FetchMode ModeID = "fetch"
	// PullMode puts the gui in pull state
	PullMode = "pull"
	// MergeMode puts the gui in merge state
	MergeMode = "merge"
	// CheckoutMode checkout selected repositories
	CheckoutMode = "checkout"

	overview Layout = 0
	focus    Layout = 1
)

var (
	mainViewFeature          = viewFeature{Name: "main", Title: " Matched Repositories "}
	mainViewFrameFeature     = viewFeature{Name: "mainframe", Title: " Matched Repositories "}
	branchViewFeature        = viewFeature{Name: "branch", Title: " Branches "}
	batchBranchViewFeature   = viewFeature{Name: "batch-branch", Title: " Select Branch "}
	suggestBranchViewFeature = viewFeature{Name: "suggest-branch", Title: " Enter New Branch Name "}
	remoteViewFeature        = viewFeature{Name: "remotes", Title: " Remotes "}
	remoteBranchViewFeature  = viewFeature{Name: "remotebranches", Title: " Remote Branches "}
	commitViewFeature        = viewFeature{Name: "commits", Title: " Commits "}
	keybindingsViewFeature   = viewFeature{Name: "keybindings", Title: " Keybindings "}
	cheatSheetViewFeature    = viewFeature{Name: "cheatsheet", Title: " Application Controls "}
	errorViewFeature         = viewFeature{Name: "error", Title: " Error "}
	dynamicViewFeature       = viewFeature{Name: "dynamic", Title: " Dynamic "}
	stashViewFeature         = viewFeature{Name: "stash", Title: " Stash "}

	fetchMode    = mode{ModeID: FetchMode, DisplayString: "Fetch", CommandString: "fetch"}
	pullMode     = mode{ModeID: PullMode, DisplayString: "Pull", CommandString: "pull"}
	mergeMode    = mode{ModeID: MergeMode, DisplayString: "Merge", CommandString: "merge"}
	checkoutMode = mode{ModeID: CheckoutMode, DisplayString: "Checkout", CommandString: "checkout"}

	modes = []mode{fetchMode, pullMode, mergeMode}
	// mainViews = []viewFeature{mainViewFeature, commitViewFeature, dynamicViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature, stashViewFeature}
	loaded = make(chan bool)
)

// New creates a Gui object and fill it's state related entities
func New(mode string, directories []string) (*Gui, error) {
	initialState := guiState{
		Directories:   directories,
		Mode:          fetchMode,
		Queue:         job.CreateJobQueue(),
		FailoverQueue: job.CreateJobQueue(),
	}
	gui := &Gui{
		State: initialState,
		mutex: &sync.Mutex{},
	}
	for _, m := range modes {
		if string(m.ModeID) == mode {
			gui.State.Mode = m
			break
		}
	}
	return gui, nil
}

// Run function runs the main loop with initial values
func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.InputEsc = true
	g.SetManagerFunc(gui.layout)

	// load repositories in background asynchronously
	go func() {
		_ = load.AsyncLoad(gui.State.Directories, gui.loadRepository, loaded)
	}()

	if err := gui.generateKeybindings(); err != nil {
		return err
	}
	if err := gui.keybindings(g); err != nil {
		return err
	}
	// mainViews = overviewViews
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

// add repository to gui's own slice and register listeners
func (gui *Gui) loadRepository(r *git.Repository) {
	rs := gui.State.Repositories
	// insertion sort implementation
	index := sort.Search(len(rs), func(i int) bool { return git.Less(r, rs[i]) })
	rs = append(rs, &git.Repository{})
	copy(rs[index+1:], rs[index:])
	rs[index] = r
	// add listener
	r.On(git.RepositoryUpdated, gui.repositoryUpdated)
	r.On(git.BranchUpdated, gui.branchUpdated)
	// update gui
	_ = gui.repositoryUpdated(nil)
	_ = gui.renderTitle()
	// take pointer back
	gui.State.Repositories = rs
	go func() {
		if <-loaded {
			v, err := gui.g.View(mainViewFrameFeature.Name)
			if err != nil {
				return
			}
			v.Title = mainViewFrameFeature.Title + fmt.Sprintf("(%d) ", len(gui.State.Repositories))
		}
	}()
}

// render title with loaded repository count
func (gui *Gui) renderTitle() error {
	v, err := gui.g.View(mainViewFrameFeature.Name)
	if err != nil {
		return err
	}
	v.Title = mainViewFrameFeature.Title + fmt.Sprintf("(%d/%d) ", len(gui.State.Repositories), len(gui.State.Directories))
	return nil
}

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) layout(g *gocui.Gui) error {
	if gui.order == overview {
		return gui.overviewLayout(g)
	} else if gui.order == focus {
		return gui.focusLayout(g)
	}
	return nil
}

// quit from the gui and end its loop
func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
