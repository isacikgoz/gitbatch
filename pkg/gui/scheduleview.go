package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
)

func (gui *Gui) updateSchedule(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View(scheduleViewFeature.Name)
    if err != nil {
        return err
    }
    out.Clear()
    if entity.Marked {
        s := "git pull " + entity.GetActiveRemote() + " " + entity.GetActiveBranch()
        fmt.Fprintln(out, s)
    } else {
        return nil
    }
    return nil
}