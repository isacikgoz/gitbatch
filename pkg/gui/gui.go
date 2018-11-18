package gui

import (
	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
	"fmt"
)

// Gui wraps the gocui Gui object which handles rendering and events
var (
	repositories []git.RepoEntity
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", 0, 0, maxX-1, int(0.65*float32(maxY))-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Matched Repositories "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		for _, r := range repositories {
			fmt.Fprintln(v, r.Name)
		}

		if _, err = setCurrentViewOnTop(g, "main"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("detail", 0, int(0.65*float32(maxY))-2, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Repository Detail "
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("keybindings", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorBlue
		v.FgColor = gocui.ColorYellow
		v.Frame = false
		fmt.Fprintln(v, "q: quit ↑ ↓: navigate space: select")
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, markRepository); err != nil {
		return err
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()

		ly := len(repositories) -1

		// if we are at the end we just return
		if cy+oy == ly {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		if err := updateDetail(g, v); err != nil {
			return err
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		if err := updateDetail(g, v); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func getSelectedRepository(g *gocui.Gui, v *gocui.View) (git.RepoEntity, error) {
	var l string
	var err error
	var r git.RepoEntity

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return r, err
	}

	for _, sr := range repositories {
		if l == sr.Name {
			return sr, nil
		}
	}
	return r, err
}

func updateDetail(g *gocui.Gui, v *gocui.View) error {
	var err error

	_, cy := v.Cursor()
	if _, err = v.Line(cy); err != nil {
		return err
	}

	out, err := g.View("detail")
	if err != nil {
		return err
	}

	out.Clear()

	if repo, err := getSelectedRepository(g, v); err != nil {
		out.Clear()
	} else {
		if list, err := repo.Repository.Remotes(); err != nil {
			return err
		} else {
			fmt.Fprintln(out, "↑" + repo.Pushables + " ↓" + repo.Pullables + " → " + repo.Branch)
			for _, r := range list {
				fmt.Fprintln(out, r)
			}
		}
	}
	return nil
}

func markRepository(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return err
	} else {
		l = l + " X"
	}

	return nil
}

// Run setup the gui with keybindings and start the mainloop
func Run(repos []git.RepoEntity) error {
	repositories = repos
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

