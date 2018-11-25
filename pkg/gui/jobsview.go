package gui

import (
    "github.com/jroimartin/gocui"
    "fmt"
    "strconv"
)

func (gui *Gui) updateJobs(g *gocui.Gui) error {
    var err error

    out, err := g.View(jobsViewFeature.Name)
    if err != nil {
        return err
    }
    out.Clear()
    pullJobs := 0
    for _, r := range gui.State.Repositories {
        if r.Marked {
            pullJobs++
        }
    }
    fcolor := white
    if pullJobs > 0 {
        fcolor = green
    }
    
    jobs := strconv.Itoa(pullJobs) + " repositories to pull"
    fmt.Fprintln(out, fcolor.Sprint(jobs))
    return nil
}