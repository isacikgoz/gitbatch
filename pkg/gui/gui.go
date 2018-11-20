package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g *gocui.Gui
	Repositories []*git.RepoEntity
}

var (
    focusedViewName string
)

// NewGui builds a new gui handler
func NewGui(entities []*git.RepoEntity) (*Gui, error) {

	gui := &Gui{
		Repositories: entities,
	}

	return gui, nil
}


func (gui *Gui) layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    if v, err := g.SetView("main", 0, 0, int(0.5*float32(maxX))-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Matched Repositories "
        v.Highlight = true
        v.SelBgColor = gocui.ColorWhite
        v.SelFgColor = gocui.ColorBlack
        v.Overwrite = true
        for _, r := range gui.Repositories {
            fmt.Fprintln(v, r.GetDisplayString())
        }

        if _, err = gui.setCurrentViewOnTop(g, "main"); err != nil {
            return err
        }
    }

    if v, err := g.SetView("status", int(0.5*float32(maxX)), 0, maxX-1, 2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Status "
        v.Wrap = false
        v.Autoscroll = false
    }

    if v, err := g.SetView("remotes", int(0.5*float32(maxX)), 3, maxX-1, int(0.25*float32(maxY))); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Remotes "
        v.Wrap = true
        v.Autoscroll = false
    }

    if v, err := g.SetView("commits", int(0.5*float32(maxX)), int(0.25*float32(maxY))+1, maxX-1, int(0.75*float32(maxY))); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Commits "
        v.Wrap = false
        v.Autoscroll = false
    }

    if v, err := g.SetView("schedule", int(0.5*float32(maxX)), int(0.75*float32(maxY))+1, maxX-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Schedule "
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
        fmt.Fprintln(v, "q: quit | ↑ ↓: navigate | space: select/deselect | a: select all | r: clear selection | enter: execute")
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

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {

    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        return err
    }
    defer g.Close()

    g.SetManagerFunc(gui.layout)

    if err := gui.keybindings(g); err != nil {
        return err
    }

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        return err
    }
    return nil
}