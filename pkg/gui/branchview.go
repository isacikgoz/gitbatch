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