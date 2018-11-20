package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    "time"
)

func (gui *Gui) openPullView(g *gocui.Gui, entity *git.RepoEntity) error {
    maxX, maxY := g.Size()
    focusedViewName = entity.Name
    v, err := g.SetView(focusedViewName, maxX/2-25, maxY/2-5, maxX/2+25, maxY/2+5)
    if err != nil {
            if err != gocui.ErrUnknownView {
                    return err
            }
            v.Title = " " + focusedViewName + " "
            v.Wrap = true
            fmt.Fprintln(v, "Pulling...")
    }
    if _, err := g.SetCurrentView(focusedViewName); err != nil {
            if err := gui.closePullView(g, entity); err != nil {
                return nil
            }
            return err
    }
    return nil
}

func (gui *Gui) closePullView(g *gocui.Gui, entity *git.RepoEntity) error {
    if g.CurrentView().Name() == entity.Name {
        if err := g.DeleteView(entity.Name); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView("main"); err != nil {
            return err
        }
    }
    return nil
}

func (gui *Gui) startPullRoutine(g *gocui.Gui, entity *git.RepoEntity) error {
    if err := gui.openPullView(g, entity); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) finalizePullRoutine(g *gocui.Gui, entity *git.RepoEntity) error {
    time.Sleep(time.Second)
    if err := gui.closePullView(g, entity); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) delView(g *gocui.Gui) error {
    if focusedViewName != "" {
        if err := g.DeleteView(focusedViewName); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView("main"); err != nil {
            return err
        }
        focusedViewName = ""
    }

    return nil
}