package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
    "strconv"
)

func (gui *Gui) updateSchedule(g *gocui.Gui) error {
    var err error

    out, err := g.View("schedule")
    if err != nil {
        return err
    }
    out.Clear()
    pullJobs := 0
    for _, r := range gui.Repositories {
        if r.Marked {
            pullJobs++
        }
    }
    jobs := strconv.Itoa(pullJobs) + " repositories to pull"
    fmt.Fprintln(out, jobs)
    return nil
}