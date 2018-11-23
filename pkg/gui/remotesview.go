package gui

import (
    "github.com/isacikgoz/gitbatch/pkg/git"
    "github.com/jroimartin/gocui"
    "fmt"
    // "strings"
)

func (gui *Gui) updateRemotes(g *gocui.Gui, entity *git.RepoEntity) (int, error) {
    var err error

    out, err := g.View("remotes")
    if err != nil {
        return 0, err
    }
    out.Clear()
    // ox, _ := out.Origin()

    currentindex := 0

    if list, err := entity.GetRemotes(); err != nil {
        return 0, err
    } else {
        for _, r := range list {
            // if strings.Contains(r, entity.Remote) {
            //     currentindex = i
            // }
            fmt.Fprintln(out, r)
        }
    }

    // if err := out.SetCursor(0, currentindex); err != nil {
            
    // if err := out.SetOrigin(0, currentindex); err != nil {
    //     return currentindex, err
    // }
    //      }

    // cx, cy := out.Cursor()
    bl := out.ViewBufferLines()
     fmt.Fprintf(out, "Bufferedlines -> %d", len(bl))
    return currentindex, nil
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

    if _, err := g.SetCurrentView("remotes"); err != nil {
        return err
    }

    if _, err = gui.updateRemotes(g, entity); err != nil {
        return err
    }

    if _, err := g.SetCurrentView("main"); err != nil {
        return err
    }

    // out, err := g.View("remotes")
    // if err != nil {
    //     return err
    // }
    // if list, err := entity.GetRemotes(); err != nil {
    //     return err
    // } else {
    //      gui.simpleCursorDown(g, out, len(list))
    // }
   
    return nil
}

// func (gui *Gui) simpleCursorDown(g *gocui.Gui, v *gocui.View, maxY int) error {
//     if v != nil {
//         cx, cy := v.Cursor()
//         ox, oy := v.Origin()

//         //ly := len(gui.State.Repositories) -1

//         // if we are at the end we just return
//         if cy+oy == maxY {
//             return nil
//         }
//         if err := v.SetCursor(cx, cy+1); err != nil {
            
//             if err := v.SetOrigin(ox, oy+1); err != nil {
//                 return err
//             }
//         }
//     }
//     return nil
// }