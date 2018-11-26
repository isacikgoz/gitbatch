package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    "strings"
)

func (gui *Gui) updateSchedule(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View(scheduleViewFeature.Name)
    if err != nil {
        return err
    }
    out.Clear()
    if entity.Marked {
        s := green.Sprint("$") + " git checkout " + entity.Branch.Name + " " + green.Sprint("✓")
        fmt.Fprintln(out, s)
        rm := entity.Remote.Reference.Name().Short()
        remote := strings.Split(rm, "/")[0]
        s = green.Sprint("$") + " git fetch " + remote
        fmt.Fprintln(out, s)
        s = green.Sprint("$") + " git merge " + entity.Remote.Name
        fmt.Fprintln(out, s)
    } else {
        return nil
    }
    return nil
}