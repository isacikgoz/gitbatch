package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateStatus(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View("status")
    if err != nil {
        return err
    }
    out.Clear()

    fmt.Fprintln(out, entity.GetStatus())

    return nil
}