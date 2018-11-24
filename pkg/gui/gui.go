package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    "time"
    "os"
    "os/exec"
    "io/ioutil"
    "log"
)

// SentinelErrors are the errors that have special meaning and need to be checked
// by calling functions. The less of these, the better
type SentinelErrors struct {
    ErrSubProcess error
}

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g *gocui.Gui
    SubProcess *exec.Cmd
    State guiState
    Errors SentinelErrors
}

type guiState struct {
    Repositories []*git.RepoEntity
    Directories  []string
}

// NewGui builds a new gui handler
func NewGui(directoies []string) (*Gui, error) {

    rs, err := git.LoadRepositoryEntities(directoies)
    if err != nil {
        return nil, err
    }
    initialState := guiState{
        Repositories: rs,
    }
	gui := &Gui{
		State: initialState,
	}

	return gui, nil
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {

    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        return err
    }
    defer g.Close()

    gui.g = g

    g.SetManagerFunc(gui.layout)

    if err := gui.keybindings(g); err != nil {
        return err
    }

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        return err
    }
    return nil
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() {
    for {
        if err := gui.Run(); err != nil {
            if err == gocui.ErrQuit {
                break
            } else if err == gui.Errors.ErrSubProcess {
                gui.SubProcess.Stdin = os.Stdin
                gui.SubProcess.Stdout = os.Stdout
                gui.SubProcess.Stderr = os.Stderr
                gui.SubProcess.Run()
                gui.SubProcess.Stdout = ioutil.Discard
                gui.SubProcess.Stderr = ioutil.Discard
                gui.SubProcess.Stdin = nil
                gui.SubProcess = nil
            } else {
                log.Fatal(err)
            }
        }
    }
}

func (gui *Gui) goEvery(g *gocui.Gui, interval time.Duration, function func(*gocui.Gui) error) {
    go func() {
        for range time.Tick(interval) {
            function(g)
        }
    }()
}

func (gui *Gui) layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    if v, err := g.SetView("main", 0, 0, int(0.55*float32(maxX))-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Matched Repositories "
        v.Highlight = true
        v.SelBgColor = gocui.ColorWhite
        v.SelFgColor = gocui.ColorBlack
        v.Overwrite = true
        for _, r := range gui.State.Repositories {
            fmt.Fprintln(v, r.DisplayString())
        }

        if _, err = gui.setCurrentViewOnTop(g, "main"); err != nil {
            return err
        }
    }

    if v, err := g.SetView("branch", int(0.55*float32(maxX)), 0, maxX-1, int(0.20*float32(maxY))-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Branches "
        v.Wrap = false
        v.Autoscroll = false
    }

    if v, err := g.SetView("remotes", int(0.55*float32(maxX)), int(0.20*float32(maxY)), maxX-1, int(0.40*float32(maxY))); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Remotes "
        v.Wrap = false
        v.Overwrite = true
    }

    if v, err := g.SetView("commits", int(0.55*float32(maxX)), int(0.40*float32(maxY))+1, maxX-1, int(0.73*float32(maxY))); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Commits "
        v.Wrap = false
        v.Autoscroll = false
    }

    if v, err := g.SetView("schedule", int(0.55*float32(maxX)), int(0.73*float32(maxY))+1, maxX-1, int(0.85*float32(maxY))); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Schedule "
        v.Wrap = true
        v.Autoscroll = true
    }

    if v, err := g.SetView("jobs", int(0.55*float32(maxX)), int(0.85*float32(maxY))+1, maxX-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Jobs "
        v.Wrap = true
        v.Autoscroll = true
    }

    if v, err := g.SetView("keybindings", -1, maxY-2, maxX, maxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.BgColor = gocui.ColorWhite
        v.FgColor = gocui.ColorBlack
        v.Frame = false
        fmt.Fprintln(v, "q: quit | ↑ ↓: navigate | space: select/deselect | c: controls | enter: execute")
    }
    return nil
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func (gui *Gui) setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
    if _, err := g.SetCurrentView(name); err != nil {
        return nil, err
    }
    return g.SetViewOnTop(name)
}

func (gui *Gui) updateKeyBindingsViewForMainView(g *gocui.Gui) error {

    v, err := g.View("keybindings")
    if err != nil {
        return err
    }

    v.Clear()
    v.BgColor = gocui.ColorWhite
    v.FgColor = gocui.ColorBlack
    v.Frame = false
    fmt.Fprintln(v, "q: quit | ↑ ↓: navigate | space: select/deselect | c: controls | enter: execute")
    return nil
}