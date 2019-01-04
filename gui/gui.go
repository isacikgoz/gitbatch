package gui

import (
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

const (
	// FetchMode puts the gui in fetch state
	FetchMode ModeID = "fetch"
	// PullMode puts the gui in pull state
	PullMode ModeID = "pull"
	// MergeMode puts the gui in merge state
	MergeMode ModeID = "merge"
)

var (
	mainViewFeature         = viewFeature{Name: "main", Title: " Matched Repositories "}
	loadingViewFeature      = viewFeature{Name: "loading", Title: " Loading in Progress "}
	branchViewFeature       = viewFeature{Name: "branch", Title: " Local Branches "}
	remoteViewFeature       = viewFeature{Name: "remotes", Title: " Remotes "}
	remoteBranchViewFeature = viewFeature{Name: "remotebranches", Title: " Remote Branches "}
	commitViewFeature       = viewFeature{Name: "commits", Title: " Commits "}
	scheduleViewFeature     = viewFeature{Name: "schedule", Title: " Schedule "}
	keybindingsViewFeature  = viewFeature{Name: "keybindings", Title: " Keybindings "}
	diffViewFeature         = viewFeature{Name: "diff", Title: " Diff Detail "}
	cheatSheetViewFeature   = viewFeature{Name: "cheatsheet", Title: " Application Controls "}
	errorViewFeature        = viewFeature{Name: "error", Title: " Error "}

	fetchMode = mode{ModeID: FetchMode, DisplayString: "Fetch", CommandString: "fetch"}
	pullMode  = mode{ModeID: PullMode, DisplayString: "Pull", CommandString: "pull"}
	mergeMode = mode{ModeID: MergeMode, DisplayString: "Merge", CommandString: "merge"}

	mainViews = []viewFeature{mainViewFeature, remoteViewFeature, remoteBranchViewFeature, branchViewFeature, commitViewFeature}
	modes     = []mode{fetchMode, pullMode, mergeMode}
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

	// If InputEsc is true, when ESC sequence is in the buffer and it doesn't
	// match any known sequence, ESC means KeyEsc.
	g.InputEsc = true
	g.SetManagerFunc(gui.layout)

	go load.AddRepositoryEntitiesAsync(gui.State.Directories, gui.addRepository)

	if err := gui.generateKeybindings(); err != nil {
		log.Error("Keybindings could not be created.")
		return err
	}
	if err := gui.keybindings(g); err != nil {
		log.Error("Keybindings could not be set.")
		return err
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Error("Error in the main loop. " + err.Error())
		return err
	}
	return nil
}

func (gui *Gui) addRepository(r *git.Repository) {
	gui.State.Repositories = append(gui.State.Repositories, r)
	r.On(git.RepositoryUpdated, gui.repositoryUpdated)
	gui.repositoryUpdated(nil)
}

// set the layout and create views with their default size, name etc. values
// TODO: window sizes can be handled better
func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(mainViewFeature.Name, 0, 0, int(0.55*float32(maxX))-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Overwrite = true
		if _, err := gui.setCurrentViewOnTop(g, mainViewFeature.Name); err != nil {
			return err
		}
	}
	if v, err := g.SetView(remoteViewFeature.Name, int(0.55*float32(maxX)), 0, maxX-1, int(0.10*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(remoteBranchViewFeature.Name, int(0.55*float32(maxX)), int(0.10*float32(maxY))+1, maxX-1, int(0.35*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = remoteBranchViewFeature.Title
		v.Wrap = false
		v.Overwrite = false
	}
	if v, err := g.SetView(branchViewFeature.Name, int(0.55*float32(maxX)), int(0.35*float32(maxY))+1, maxX-1, int(0.60*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = branchViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(commitViewFeature.Name, int(0.55*float32(maxX)), int(0.60*float32(maxY))+1, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitViewFeature.Title
		v.Wrap = false
		v.Autoscroll = false
	}
	if v, err := g.SetView(keybindingsViewFeature.Name, -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		v.Frame = false
		gui.updateKeyBindingsView(g, mainViewFeature.Name)
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
