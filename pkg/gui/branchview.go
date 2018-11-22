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
    branches, err := entity.GetBranches()
    if err != nil {
        return err
    }
    for _, b := range branches {
        fmt.Fprintln(out, b)
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