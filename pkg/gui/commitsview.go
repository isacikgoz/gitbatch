package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) updateCommits(g *gocui.Gui, entity *git.RepoEntity) error {
	var err error
	out, err := g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	out.Clear()

	currentindex := 0
	totalcommits := len(entity.Commits)
	for i, c := range entity.Commits {
		if c.Hash == entity.Commit.Hash {
			currentindex = i
			fmt.Fprintln(out, selectionIndicator+green.Sprint(c.Hash[:git.Hashlimit])+" "+c.Message)
			continue
		}
		fmt.Fprintln(out, tab+cyan.Sprint(c.Hash[:git.Hashlimit])+" "+c.Message)
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
	v, err := g.SetView(commitdetailViewFeature.Name, 5, 3, maxX-5, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = commitdetailViewFeature.Title
		v.Overwrite = true
		v.Wrap = true

		main, _ := g.View(mainViewFeature.Name)

		entity, err := gui.getSelectedRepository(g, main)
		if err != nil {
			return err
		}
		commit := entity.Commit
		commitDetail := "Hash: " + cyan.Sprint(commit.Hash) + "\n" + "Author: " + commit.Author + "\n" + commit.Time.String() + "\n" + "\n" + "\t\t" + commit.Message + "\n"
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

	gui.updateKeyBindingsView(g, commitdetailViewFeature.Name)
	if _, err := g.SetCurrentView(commitdetailViewFeature.Name); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) closeCommitDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

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
