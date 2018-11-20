package gui

import (
    "fmt"
    "sync"
    "github.com/isacikgoz/gitbatch/pkg/utils"
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
)

func (gui *Gui) refreshViews(g *gocui.Gui, entity *git.RepoEntity) error {

    if err := gui.updateRemotes(g, entity); err != nil {
        return err
    }

    if err := gui.updateStatus(g, entity); err != nil {
        return err
    }

    if err := gui.updateCommits(g, entity); err != nil {
         return err
    }

    if err := gui.updateSchedule(g); err != nil {
         return err
    }

    return nil
}

func (gui *Gui) cursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        ox, oy := v.Origin()

        ly := len(gui.Repositories) -1

        // if we are at the end we just return
        if cy+oy == ly {
            return nil
        }
        if err := v.SetCursor(cx, cy+1); err != nil {
            
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
        if entity, err := gui.getSelectedRepository(g, v); err != nil {
            return err
        } else {
            gui.refreshViews(g, entity)
        }
    }
    return nil
}

func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                return err
            }
        }
        if entity, err := gui.getSelectedRepository(g, v); err != nil {
            return err
        } else {
            gui.refreshViews(g, entity)
        }
    }
    return nil
}

func (gui *Gui) getSelectedRepository(g *gocui.Gui, v *gocui.View) (*git.RepoEntity, error) {
    var l string
    var err error
    var r *git.RepoEntity

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        return r, err
    }

    for _, sr := range gui.Repositories {
        if l == sr.Name {
            return sr, nil
        }
    }
    return r, err
}

func (gui *Gui) markRepository(g *gocui.Gui, v *gocui.View) error {

    if r, err := gui.getSelectedRepository(g, v); err != nil {
        return err
    } else {
        if err != nil {
            return err
        }
        if r.Marked != true {
            r.Mark()
        } else {
            r.UnMark()
        }
        gui.refreshMain(g)
        gui.updateSchedule(g)
    }
    
    return nil
}

func (gui *Gui) markAllRepositories(g *gocui.Gui, v *gocui.View) error {
    for _, r := range gui.Repositories {
        r.Mark()
    }
    if err := gui.refreshMain(g); err !=nil {
        return err
    }
    gui.updateSchedule(g)
    return nil
}

func (gui *Gui) unMarkAllRepositories(g *gocui.Gui, v *gocui.View) error {
    for _, r := range gui.Repositories {
        r.UnMark()
    }
    if err := gui.refreshMain(g); err !=nil {
        return err
    }
    gui.updateSchedule(g)
    return nil
}

func (gui *Gui) refreshMain(g *gocui.Gui) error {
    
    mainView, err := g.View("main")
    if err != nil {
        return err
    }
    mainView.Clear()
    for _, r := range gui.Repositories {
        fmt.Fprintln(mainView, r.GetDisplayString())
    }
    return nil
}

// if the cursor down past the last item, move it to the last line
func (gui *Gui) correctCursor(v *gocui.View) error {
    cx, cy := v.Cursor()
    ox, oy := v.Origin()
    width, height := v.Size()
    maxY := height - 1
    ly := width - 1
    if oy+cy <= ly {
        return nil
    }
    newCy := utils.Min(ly, maxY)
    if err := v.SetCursor(cx, newCy); err != nil {
        return err
    }
    if err := v.SetOrigin(ox, ly-newCy); err != nil {
        return err
    }
    return nil
}

func (gui *Gui) getMarkedEntities() (rs []*git.RepoEntity, err error) {
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, r := range gui.Repositories {
        wg.Add(1)
        go func(repo *git.RepoEntity){
            defer wg.Done()
            if repo.Marked {
                mu.Lock()
                rs = append(rs, repo)
                mu.Unlock()
            }
        }(r)
    }
    wg.Wait()

    return rs, nil
}