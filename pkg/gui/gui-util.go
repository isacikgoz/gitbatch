package gui

import (
	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/isacikgoz/gitbatch/pkg/utils"
	"github.com/jroimartin/gocui"
)

var (
	blue    = color.New(color.FgBlue)
	green   = color.New(color.FgGreen)
	red     = color.New(color.FgRed)
	cyan    = color.New(color.FgCyan)
	yellow  = color.New(color.FgYellow)
	white   = color.New(color.FgWhite)
	magenta = color.New(color.FgMagenta)
)

func (gui *Gui) refreshViews(g *gocui.Gui, entity *git.RepoEntity) error {

	if err := gui.updateRemotes(g, entity); err != nil {
		return err
	}
	if err := gui.updateBranch(g, entity); err != nil {
		return err
	}
	if err := gui.updateRemoteBranches(g, entity); err != nil {
		return err
	}
	if err := gui.updateCommits(g, entity); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) switchMode(g *gocui.Gui, v *gocui.View) error {
	switch mode := gui.State.Mode.ModeID; mode {
	case FetchMode:
		gui.State.Mode = pullMode
	case PullMode:
		gui.State.Mode = mergeMode
	case MergeMode:
		gui.State.Mode = fetchMode
	default:
		gui.State.Mode = fetchMode
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

func (gui *Gui) setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
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

func (gui *Gui) smartAnchorRelativeToLine(v *gocui.View, currentindex, totallines int) error {

	_, y := v.Size()
	if currentindex >= int(0.5*float32(y)) && totallines-currentindex+int(0.5*float32(y)) >= y {
		if err := v.SetOrigin(0, currentindex-int(0.5*float32(y))); err != nil {
			return err
		}
	} else if totallines-currentindex < y && totallines > y {
		if err := v.SetOrigin(0, totallines-y); err != nil {
			return err
		}
	} else if totallines-currentindex <= int(0.5*float32(y)) && totallines > y-1 && currentindex > y {
		if err := v.SetOrigin(0, currentindex-int(0.5*float32(y))); err != nil {
			return err
		}
	} else {
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
	}
	return nil
}

func writeRightHandSide(v *gocui.View, text string, cx, cy int) error {
	runes := []rune(text)
	tl := len(runes)
	lx, _ := v.Size()
	v.MoveCursor(lx-tl, cy-1, true)
	for i := tl - 1; i >= 0; i-- {
		v.EditDelete(true)
		v.EditWrite(runes[i])
	}
	v.SetCursor(cx, cy)
	return nil
}
