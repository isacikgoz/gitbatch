package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateBranch(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View("branch")
    if err != nil {
        return err
    }
    out.Clear()

    currentindex := 0
    totalbranches := 0
    if branches, err := entity.LocalBranches(); err != nil {
        return err
    } else {
        totalbranches = len(branches)
        for i, b := range branches {
            if b == entity.Branch {
                currentindex = i
                fmt.Fprintln(out, selectionIndicator() + b)
                continue
            } 
            fmt.Fprintln(out, tab() + b)
        }
    }
    if err = gui.smartAnchorRelativeToLine(out, currentindex, totalbranches); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) nextBranch(g *gocui.Gui, v *gocui.View) error {
    var err error

    entity, err := gui.getSelectedRepository(g, v)
    if err != nil {
        return err
    }

    if err = entity.Checkout(entity.NextBranch()); err != nil {
        return err
    }

    if err = gui.updateBranch(g, entity); err != nil {
        return err
    }

    if err = gui.updateCommits(g, entity); err != nil {
        return err
    }

    if err = gui.refreshMain(g); err != nil {
        return err
    }
    return nil
}