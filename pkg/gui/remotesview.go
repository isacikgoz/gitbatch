package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateRemotes(g *gocui.Gui, entity *git.RepoEntity) error {
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