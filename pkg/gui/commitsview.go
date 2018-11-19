package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateCommits(g *gocui.Gui, entity git.RepoEntity) error {
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