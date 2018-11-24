package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    "log"
)

func (gui *Gui) updateCommits(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View("commits")
    if err != nil {
        return err
    }
    out.Clear()

    currentindex := 0
    if commits, err := entity.Commits(); err != nil {
        return err
    } else {
        for i, c := range commits {
            if c[:git.Hashlimit] == entity.Commit {
                currentindex = i
                fmt.Fprintln(out, selectionIndicator() + c)
                continue
            } 
            fmt.Fprintln(out, tab() + c)
        }
    }
    _, y := out.Size()
    if currentindex > y-1 {
        if err := out.SetOrigin(0, currentindex - int(0.5*float32(y))); err != nil {
            return err
        }
    } else {
        if err := out.SetOrigin(0, 0); err != nil {
            return err
        }
    }
    return nil
}

func (gui *Gui) nextCommit(g *gocui.Gui, v *gocui.View) error {
    var err error

    entity, err := gui.getSelectedRepository(g, v)
    if err != nil {
        return err
    }

    if err = entity.NextCommit(); err != nil {
        return err
    }

    if err = gui.updateCommits(g, entity); err != nil {
        return err
    }

    return nil
}

func (gui *Gui) showCommitDetail(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()

    v, err := g.SetView("commitdetail", maxX/2-35, maxY/2-5, maxX/2+35, maxY/2+5)
    if err != nil {
        if err != gocui.ErrUnknownView {
             return err
        }
        v.Title = " Commit Detail "
        v.Highlight = true
        v.Overwrite = true

        main, _ := g.View("main")

        entity, err := gui.getSelectedRepository(g, main)
        if err != nil {
            log.Fatal(err)
            return err
        }

        detail, err := entity.CommitDetail()
        if err != nil {
            return err
        }
        fmt.Fprintln(v, detail)
        diff, err := entity.Diff(entity.Commit)
        if err != nil {
            return err
        }
        fmt.Fprintln(v, diff)
    }
    
    gui.updateKeyBindingsViewForCommitDetailView(g)
    if _, err := g.SetCurrentView("commitdetail"); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) closeCommitDetailView(g *gocui.Gui, v *gocui.View) error {

        if err := g.DeleteView(v.Name()); err != nil {
            return nil
        }
        if _, err := g.SetCurrentView("main"); err != nil {
            return err
        }
        gui.updateKeyBindingsViewForMainView(g)

    return nil
}

func (gui *Gui) updateKeyBindingsViewForCommitDetailView(g *gocui.Gui) error {

    v, err := g.View("keybindings")
    if err != nil {
        return err
    }
    v.Clear()
    v.BgColor = gocui.ColorWhite
    v.FgColor = gocui.ColorBlack
    v.Frame = false
    fmt.Fprintln(v, "c: cancel | ↑ ↓: navigate")
    return nil
}