package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

// open diff view for the selcted commit
func (gui *Gui) openCommitDiffView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(commitDiffViewFeature.Name, 5, 3, maxX-5, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitDiffViewFeature.Title
		v.Overwrite = true
		v.Wrap = true

		main, _ := g.View(mainViewFeature.Name)

		entity, err := gui.getSelectedRepository(g, main)
		if err != nil {
			return err
		}
		commit := entity.Commit
		commitDetail := "Hash: " + cyan.Sprint(commit.Hash) + "\n" + "Author: " + commit.Author +
		 "\n" + commit.Time + "\n" + "\n" + "\t\t" + commit.Message + "\n"
		fmt.Fprintln(v, commitDetail)
		diff, err := entity.Diff(entity.Commit.Hash)
		if err != nil {
			return err
		}
		colorized := colorizeDiff(diff)
		for _, line := range colorized {
			fmt.Fprintln(v, line)
		}
	}

	gui.updateKeyBindingsView(g, commitDiffViewFeature.Name)
	if _, err := g.SetCurrentView(commitDiffViewFeature.Name); err != nil {
		return err
	}
	return nil
}

// close the opened diff view
func (gui *Gui) closeCommitDiffView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

// cursor down acts like half-page down for faster scrolling
func (gui *Gui) commitCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, vy := v.Size()

		// TODO: do something when it hits bottom
		if err := v.SetOrigin(ox, oy+vy/2); err != nil {
			return err
		}
	}
	return nil
}

// cursor up acts like half-page up for faster scrolling
func (gui *Gui) commitCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, vy := v.Size()

		if oy-vy/2 > 0 {
			if err := v.SetOrigin(ox, oy-vy/2); err != nil {
				return err
			}
		} else if oy-vy/2 <= 0 {
			if err := v.SetOrigin(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

// colorize the plain diff text collected from system output
// the style is near to original diff command
func colorizeDiff(original string) (colorized []string) {
	colorized = strings.Split(original, "\n")
	re := regexp.MustCompile(`@@ .+ @@`)
	for i, line := range colorized {
		if len(line) > 0 {
			if line[0] == '-' {
				colorized[i] = red.Sprint(line)
			} else if line[0] == '+' {
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
