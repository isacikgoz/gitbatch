package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/queue"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// Gui struct hold the gocui struct along with the gui's state, also keybindings
// are tied with this struct in order to render those in different occasions
type Gui struct {
	g           *gocui.Gui
	KeyBindings []*KeyBinding
	State       guiState
}

// guiState struct holds the repositories, directiories, mode and queue of the
// gui object. These values are not static
type guiState struct {
	Repositories []*git.RepoEntity
	Directories  []string
	Mode         mode
	Queue        *queue.JobQueue
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
	commitDiffViewFeature   = viewFeature{Name: "commitdiff", Title: " Commit Detail "}
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
		Directories: directoies,
		Mode:        fetchMode,
		Queue:       queue.CreateJobQueue(),
	}
	gui := &Gui{
		State: initialState,
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

	// start an async view apart from this loop to show loading screen
	go func(g_ui *Gui) {
		maxX, maxY := g.Size()
		// TODO: view size can be handled in a more smart way
		v, err := g.SetView(loadingViewFeature.Name, maxX/2-10, maxY/2-1, maxX/2+10, maxY/2+1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				log.Warn("Loading view cannot be created.")
				return
			}
			fmt.Fprintln(v, "Loading...")
		}
		if _, err := g.SetCurrentView(loadingViewFeature.Name); err != nil {
			log.Warn("Loading view cannot be focused.")
			return
		}
		rs, err := git.LoadRepositoryEntities(g_ui.State.Directories)
		if err != nil {
			g.Close()
			log.Fatal(err)
			return
		}
		g_ui.State.Repositories = rs
		gui.fillMain(g)
	}(gui)

	defer g.Close()
	gui.g = g
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(gui.layout)

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
	err := gui.nextViewOfGroup(g, v, mainViews)
	return err
}

// focus to previous view
func (gui *Gui) previousMainView(g *gocui.Gui, v *gocui.View) error {
	err := gui.previousViewOfGroup(g, v, mainViews)
	return err
}
// quit from the gui and end its loop
func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
