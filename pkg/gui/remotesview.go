package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateRemotes(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View(remoteViewFeature.Name)
    if err != nil {
        return err
    }
    out.Clear()

    currentindex := 0
    totalRemotes := len(entity.Remotes)
    for i, r := range entity.Remotes {
        if r.Name == entity.Remote.Name {
            currentindex = i
            fmt.Fprintln(out, selectionIndicator() + r.Name)
            continue
        } 
        fmt.Fprintln(out, tab() + r.Name)
    } 
    if err = gui.smartAnchorRelativeToLine(out, currentindex, totalRemotes); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) nextRemote(g *gocui.Gui, v *gocui.View) error {
    var err error
    entity, err := gui.getSelectedRepository(g, v)
    if err != nil {
        return err
    }

    if err = entity.NextRemote(); err != nil {
        return err
    }

    if err = gui.updateRemotes(g, entity); err != nil {
        return err
    }
    return nil
}