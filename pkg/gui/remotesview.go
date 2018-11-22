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

func (gui *Gui) nextRemote(g *gocui.Gui, v *gocui.View) error {
    var err error

    entity, err := gui.getSelectedRepository(g, v)
    if err != nil {
        return err
    }
    if _, err = entity.NextRemote(); err != nil {
        return err
    }

    if err = gui.updateRemotes(g, entity); err != nil {
        return err
    }

    return nil
}