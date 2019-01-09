package gui

import (
	"fmt"
	"sort"
	"sync"

	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/isacikgoz/gitbatch/core/job"
	"github.com/isacikgoz/gitbatch/core/load"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
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

// guiState struct holds the repositories, directiories, mode and queue of the
// gui object. These values are not static
type guiState struct {
	Repositories  []*git.Repository
	Directories   []string
	Mode          mode
	Queue         *job.JobQueue
	FailoverQueue *job.JobQueue
}

// this struct encapsulates the name and title of a view. the name of a view is
// passed around so much it is added so that I don't need to wirte names again
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

// ModeID is the mode indicator for the gui
type ModeID string

type Layout int

const (
	// FetchMode puts the gui in fetch state
	FetchMode ModeID = "fetch"
	// PullMode puts the gui in pull state
	PullMode ModeID = "pull"
	// MergeMode puts the gui in merge state
	MergeMode ModeID = "merge"

	overview Layout = 0

	focus Layout = 1
)

var (
	mainViewFeature         = viewFeature{Name: "main", Title: " Matched Repositories "}
	loadingViewFeature      = viewFeature{Name: "loading", Title: " Loading in Progress "}
	branchViewFeature       = viewFeature{Name: "branch", Title: " Branches "}
	remoteViewFeature       = viewFeature{Name: "remotes", Title: " Remotes "}
	remoteBranchViewFeature = viewFeature{Name: "remotebranches", Title: " Remote Branches "}
	commitViewFeature       = viewFeature{Name: "commits", Title: " Commits "}
	scheduleViewFeature     = viewFeature{Name: "schedule", Title: " Schedule "}
	keybindingsViewFeature  = viewFeature{Name: "keybindings", Title: " Keybindings "}
	diffViewFeature         = viewFeature{Name: "diff", Title: " Diff Detail "}
	cheatSheetViewFeature   = viewFeature{Name: "cheatsheet", Title: " Application Controls "}
	errorViewFeature        = viewFeature{Name: "error", Title: " Error "}
	detailViewFeature       = viewFeature{Name: "details", Title: " Details "}
	stashedViewFeature      = viewFeature{Name: "stashed", Title: " Stash "}

	fetchMode = mode{ModeID: FetchMode, DisplayString: "Fetch", CommandString: "fetch"}
	pullMode  = mode{ModeID: PullMode, DisplayString: "Pull", CommandString: "pull"}
	mergeMode = mode{ModeID: MergeMode, DisplayString: "Merge", CommandString: "merge"}

	modes     = []mode{fetchMode, pullMode, mergeMode}
	mainViews = []viewFeature{mainViewFeature, commitViewFeature, detailViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature}
	loaded    = make(chan bool)
)

// NewGui creates a Gui opject and fill it's state related entites
func NewGui(mode string, directoies []string) (*Gui, error) {
	initialState := guiState{
		Directories:   directoies,
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
	go load.AsyncLoad(gui.State.Directories, gui.loadRepository, loaded)

	if err := gui.generateKeybindings(); err != nil {
		log.Error("Keybindings could not be created.")
		return err
	}
	if err := gui.keybindings(g); err != nil {
		log.Error("Keybindings could not be set.")
		return err
	}
	mainViews = overviewViews
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Error("Error in the main loop. " + err.Error())
		return err
	}
	return nil
}

func (gui *Gui) loadRepository(r *git.Repository) {
	rs := gui.State.Repositories

	// insertion sort implementation
	index := sort.Search(len(rs), func(i int) bool { return git.Less(r, rs[i]) })
	rs = append(rs, &git.Repository{})
	copy(rs[index+1:], rs[index:])
	rs[index] = r
	// add listener
	r.On(git.RepositoryUpdated, gui.repositoryUpdated)
	// update gui
	gui.repositoryUpdated(nil)
	gui.renderTitle()
	// take pointer back
	gui.State.Repositories = rs
	go func() {
		if <-loaded {
			v, err := gui.g.View(mainViewFeature.Name)
			if err != nil {
				log.Warn(err.Error())
				return
			}
			v.Title = mainViewFeature.Title + fmt.Sprintf("(%d) ", len(gui.State.Repositories))
		}
	}()
}

func (gui *Gui) renderTitle() error {
	v, err := gui.g.View(mainViewFeature.Name)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	v.Title = mainViewFeature.Title + fmt.Sprintf("(%d/%d) ", len(gui.State.Repositories), len(gui.State.Directories))
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

// focus to next view
func (gui *Gui) nextMainView(g *gocui.Gui, v *gocui.View) error {
	return gui.nextViewOfGroup(g, v, mainViews)
}

// focus to previous view
func (gui *Gui) previousMainView(g *gocui.Gui, v *gocui.View) error {
	return gui.previousViewOfGroup(g, v, mainViews)
}

// quit from the gui and end its loop
func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
