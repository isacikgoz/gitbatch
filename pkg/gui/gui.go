package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

// Gui wraps the gocui Gui object which handles rendering and events
var (
    repositories []git.RepoEntity
)

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    if v, err := g.SetView("main", 0, 0, int(0.5*float32(maxX))-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = " Matched Repositories "
        v.Highlight = true
        v.SelBgColor = gocui.ColorWhite
        v.SelFgColor = gocui.ColorBlack

        for _, r := range repositories {
            fmt.Fprintln(v, r.Name)
        }

        if _, err = setCurrentViewOnTop(g, "main"); err != nil {
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
        v.Autoscroll = true
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
        fmt.Fprintln(v, "q: quit ↑ ↓: navigate space: select")
    }
    return nil
}

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, markRepository); err != nil {
        return err
    }
    return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        ox, oy := v.Origin()

        ly := len(repositories) -1

        // if we are at the end we just return
        if cy+oy == ly {
            return nil
        }
        if err := v.SetCursor(cx, cy+1); err != nil {
            
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
        if entity, err := getSelectedRepository(g, v); err != nil {
            return err
        } else {
            if err := updateRemotes(g, entity); err != nil {
                return err
            }

            if err := updateStatus(g, entity); err != nil {
                return err
            }

            if err := updateCommits(g, entity); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                return err
            }
        }
        if entity, err := getSelectedRepository(g, v); err != nil {
            return err
        } else {
            if err := updateRemotes(g, entity); err != nil {
                return err
            }

            if err := updateStatus(g, entity); err != nil {
                return err
            }

            if err := updateCommits(g, entity); err != nil {
                return err
            }
        }
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
    if _, err := g.SetCurrentView(name); err != nil {
        return nil, err
    }
    return g.SetViewOnTop(name)
}

func getSelectedRepository(g *gocui.Gui, v *gocui.View) (git.RepoEntity, error) {
    var l string
    var err error
    var r git.RepoEntity

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        return r, err
    }

    for _, sr := range repositories {
        if l == sr.Name {
            return sr, nil
        }
    }
    return r, err
}

func updateRemotes(g *gocui.Gui, entity git.RepoEntity) error {
    var err error

    out, err := g.View("remotes")
    if err != nil {
        return err
    }
    out.Clear()

    if list, err := entity.GetRemotes(); err != nil {
        return err
    } else {
        for _, r := range list {
            fmt.Fprintln(out, r)
        }
    }

    return nil
}

func updateStatus(g *gocui.Gui, entity git.RepoEntity) error {
    var err error

    out, err := g.View("status")
    if err != nil {
        return err
    }
    out.Clear()

    fmt.Fprintln(out, entity.GetStatus())

    return nil
}


func updateCommits(g *gocui.Gui, entity git.RepoEntity) error {
    var err error

    out, err := g.View("commits")
    if err != nil {
        return err
    }
    out.Clear()

    if list, err := entity.GetCommits(); err != nil {
        return err
    } else {
        for _, c := range list {
            fmt.Fprintln(out, c)
        }
    }

    return nil
}

func markRepository(g *gocui.Gui, v *gocui.View) error {
    var l string
    var err error

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        return err
    } else {
        l = l + " X"
    }

    return nil
}

// Run setup the gui with keybindings and start the mainloop
func Run(repos []git.RepoEntity) error {
    repositories = repos
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        return err
    }
    defer g.Close()

    g.SetManagerFunc(layout)

    if err := keybindings(g); err != nil {
        return err
    }

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        return err
    }
    return nil
}

