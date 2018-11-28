package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/job"
	"github.com/jroimartin/gocui"
)

type Gui struct {
	g           *gocui.Gui
	KeyBindings []*KeyBinding
	State       guiState
}

type guiState struct {
	Repositories []*git.RepoEntity
	Directories  []string
	Mode         mode
	Queue        *job.JobQueue
}

type viewFeature struct {
	Name  string
	Title string
}

type mode struct {
	ModeID        ModeID
	DisplayString string
	CommandString string
	ExecString    string
}

type ModeID int8

const (
	FetchMode ModeID = 0
	PullMode  ModeID = 1
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
	commitdetailViewFeature = viewFeature{Name: "commitdetail", Title: " Commit Detail "}
	cheatSheetViewFeature   = viewFeature{Name: "cheatsheet", Title: " Application Controls "}
	errorViewFeature        = viewFeature{Name: "error", Title: " Error "}

	fetchMode = mode{ModeID: FetchMode, DisplayString: "Fetch", CommandString: "fetch", ExecString: "repositories will be fetched"}
	pullMode  = mode{ModeID: PullMode, DisplayString: "Pull", CommandString: "pull", ExecString: "repositories will be pulled"}
)

func NewGui(directoies []string) (*Gui, error) {
	initialState := guiState{
		Directories: directoies,
		Mode:        fetchMode,
		Queue:       job.CreateJobQueue(),
	}
	gui := &Gui{
		State: initialState,
	}
	return gui, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	go func(g_ui *Gui) {
		maxX, maxY := g.Size()
		v, err := g.SetView(loadingViewFeature.Name, maxX/2-10, maxY/2-1, maxX/2+10, maxY/2+1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return
			}
			fmt.Fprintln(v, "Loading...")
		}
		if _, err := g.SetCurrentView(loadingViewFeature.Name); err != nil {
			return
		}
		rs, err := git.LoadRepositoryEntities(g_ui.State.Directories)
		if err != nil {
			return
		}
		g_ui.State.Repositories = rs
		gui.fillMain(g)
	}(gui)

	defer g.Close()
	gui.g = g
	g.SetManagerFunc(gui.layout)

	if err := gui.generateKeybindings(); err != nil {
		return err
	}
	if err := gui.keybindings(g); err != nil {
		return err
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(mainViewFeature.Name, 0, 0, int(0.55*float32(maxX))-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = mainViewFeature.Title
		v.Highlight = true
		v.SelFgColor = gocui.ColorGreen
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

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
