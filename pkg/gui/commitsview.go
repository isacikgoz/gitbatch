package gui

import (
    "github.com/fatih/color"
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    "strings"
    "regexp"
)

func (gui *Gui) updateCommits(g *gocui.Gui, entity *git.RepoEntity) error {
    var err error

    out, err := g.View("commits")
    if err != nil {
        return err
    }
    out.Clear()

    cyan := color.New(color.FgCyan)
    green := color.New(color.FgGreen)
    totalcommits := 0
    currentindex := 0
    if commits, err := entity.Commits(); err != nil {
        return err
    } else {
        totalcommits = len(commits)
        for i, c := range commits {
            if c.Hash == entity.Commit.Hash {
                currentindex = i
                fmt.Fprintln(out, selectionIndicator() + green.Sprint(c.Hash[:git.Hashlimit]) + " " + c.Message)
                continue
            } 
            fmt.Fprintln(out, tab() + cyan.Sprint(c.Hash[:git.Hashlimit]) + " " + c.Message)
        }
    }
    if err = gui.smartAnchorRelativeToLine(out, currentindex, totalcommits); err != nil {
        return err
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
    cyan := color.New(color.FgCyan)
    v, err := g.SetView("commitdetail", 5, 3, maxX-5, maxY-3)
    if err != nil {
        if err != gocui.ErrUnknownView {
             return err
        }
        v.Title = " Commit Detail "
        v.Overwrite = true
        v.Wrap = true

        main, _ := g.View("main")

        entity, err := gui.getSelectedRepository(g, main)
        if err != nil {
            return err
        }
        commit := entity.Commit
        commitDetail := "Hash: " + cyan.Sprint(commit.Hash) + "\n" + "Author: " + commit.Author + "\n" + commit.Time.String() + "\n" + "\n" + "\t" + commit.Message + "\n"
        fmt.Fprintln(v, commitDetail)
        diff, err := entity.Diff(entity.Commit.Hash)
        if err != nil {
            return err
        }
        colorized := colorizeDiff(diff)
        for _, line := range colorized{
            fmt.Fprintln(v, line)
        }
        
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

func (gui *Gui) commitCursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        ox, oy := v.Origin()

        // TODO: do something when it hits bottom
        // also it woulf be nice to emulate less buttons
        // or at least page-up/page-down buttons

        // ly := len(gui.State.Repositories) -1

        // // if we are at the end we just return
        // if cy+oy == ly {
        //     return nil
        // }
        if err := v.SetCursor(cx, cy+1); err != nil {
            
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
    }
    return nil
}

func (gui *Gui) commitCursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                return err
            }
        }
    }
    return nil
}

func colorizeDiff(original string) (colorized []string) {
    cyan := color.New(color.FgCyan)
    green := color.New(color.FgGreen)
    red := color.New(color.FgRed)
    colorized = strings.Split(original, "\n")
    re := regexp.MustCompile(`@@ .+ @@`)
    for i, line := range colorized {
        if len(line) > 0 {
            if line[0] == '-' {
                colorized[i] = red.Sprint(line)
            } else if line [0] == '+' {
                colorized[i] = green.Sprint(line)
            } else if re.MatchString(line) {
                s := re.FindString(line)
                colorized[i] = cyan.Sprint(s) + line[len(s):]
            } else {
                continue
            }
        } else {
            continue
        }
        
    }
    return colorized
}