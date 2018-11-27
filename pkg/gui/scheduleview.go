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
        
        rm := entity.Remote.Reference.Name().Short()
        switch mode := gui.State.Mode.ModeID; mode {
        case FetchMode:
            remote := strings.Split(rm, "/")[0]
            s := green.Sprint("$") + " git fetch " + remote
            fmt.Fprintln(out, s)
        case PullMode:
            s := green.Sprint("$") + " git checkout " + entity.Branch.Name + " " + green.Sprint("✓")
            fmt.Fprintln(out, s)
            remote := strings.Split(rm, "/")[0]
            s = green.Sprint("$") + " git fetch " + remote
            fmt.Fprintln(out, s)
            s = green.Sprint("$") + " git merge " + entity.Remote.Name
            fmt.Fprintln(out, s)
        default:
            fmt.Fprintln(out, "No mode selected")
        }
    } else {
        return nil
    }
    return nil
}