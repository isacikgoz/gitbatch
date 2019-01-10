package gui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) focusToRepository(g *gocui.Gui, v *gocui.View) error {
	mainViews = focusViews
	r := gui.getSelectedRepository()
	gui.order = focus

	if _, err := g.SetCurrentView(commitViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, commitViewFeature.Name)

	r.State.Branch.InitializeCommits(r)

	gui.renderCommits(r)
	gui.initFocusStat(r)
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.renderMain()
	})
	return nil
}

func (gui *Gui) focusBackToMain(g *gocui.Gui, v *gocui.View) error {
	mainViews = overviewViews
	gui.order = overview

	if _, err := g.SetCurrentView(mainViewFeature.Name); err != nil {
		return err
	}
	gui.updateKeyBindingsView(g, mainViewFeature.Name)
	return nil
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) commitCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		ly := len(v.BufferLines()) - 1

		// if we are at the end we just return
		if cy+oy == ly-1 {
			return nil
		}
		v.EditDelete(true)
		pos := cy + oy + 1
		adjustAnchor(pos, ly, v)
		gui.commitDetail(pos)
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) commitCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ly := len(v.BufferLines()) - 1
		v.EditDelete(true)
		pos := cy + oy - 1
		adjustAnchor(pos, ly, v)
		if pos >= 0 {
			gui.commitDetail(cy + oy - 1)
		}
	}
	return nil
}

// updates the commitsview for given entity
func (gui *Gui) renderCommits(r *git.Repository) error {
	v, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	cs := r.State.Branch.Commits
	bc := r.State.Branch.State.Commit
	si := 0
	// if r.State.Branch.Clean {
	// 	fmt.Fprintln(v, " "+cyan.Sprint("*******")+" "+green.Sprint("Current State"))
	// } else {
	fmt.Fprintln(v, " "+yellow.Sprint("*******")+" "+yellow.Sprint("Current State"))
	// }

	for i, c := range cs {
		if c.Hash == bc.Hash {
			si = i
			fmt.Fprintln(v, ws+commitLabel(c, false))
			continue
		}

		fmt.Fprintln(v, tab+commitLabel(c, false))
	}
	adjustAnchor(si, len(cs), v)
	return nil
}

func (gui *Gui) commitDetail(ix int) error {

	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	r := gui.getSelectedRepository()
	v.Clear()
	if ix == 0 {
		return gui.initFocusStat(r)
	}
	v.Title = detailViewFeature.Title
	c := r.State.Branch.Commits[ix-1]
	done := make(chan bool)
	var stat string
	go func() {
		stat = c.DiffStat(done)
	}()
	ld := "loading stats..."
	fmt.Fprintf(v, "%s\n", decorateCommit(c.String()))
	fmt.Fprintf(v, "%s\n", red.Sprint(ld))

	go func(gui *Gui) {
		if <-done {
			if !strings.Contains(strings.Join(v.BufferLines(), c.Hash), ld) {
				return
			}
			v.Clear()
			fmt.Fprintf(v, "%s\n", decorateCommit(c.String()))
			gui.g.Update(func(g *gocui.Gui) error {
				v, err := gui.g.View(detailViewFeature.Name)
				if err != nil {
					return err
				}
				fmt.Fprintf(v, decorateDiffStat(stat, true))
				return nil
			})
		}
	}(gui)

	return nil
}

func (gui *Gui) commitDiff(g *gocui.Gui, v *gocui.View) error {
	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}
	vcm, err := gui.g.View(commitViewFeature.Name)
	if err != nil {
		return err
	}
	_, oy := vcm.Origin()
	_, cy := vcm.Cursor()
	ix := oy + cy

	r := gui.getSelectedRepository()
	if ix == 0 {
		p, err := command.PlainDiff(r)
		if err != nil {
			return err
		}
		var s string
		for _, d := range colorizeDiff(p) {
			s = s + "\n" + d
		}
		if len(s) <= 1 {
			return nil
		}

		v.Clear()
		fmt.Fprintf(v, s)
		return nil
	}

	v.Clear()
	c := r.State.Branch.Commits[ix-1]
	if ix+1 > len(r.State.Branch.Commits) {
		ix = ix - 1
	}
	p, err := r.State.Branch.Commits[ix].C.Patch(c.C)
	if err != nil {
		return err
	}
	var s string
	for _, d := range colorizeDiff(p.String()) {
		s = s + "\n" + d
	}
	fmt.Fprintf(v, s)
	return nil
}

// moves cursor down for a page size
func (gui *Gui) dpageDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, _ := v.Cursor()
		_, vy := v.Size()
		lr := len(v.BufferLines())
		if lr < vy {
			return nil
		}
		if oy+vy >= lr-vy {
			if err := v.SetOrigin(ox, lr-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy+vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return nil
}

// moves cursor up for a page size
func (gui *Gui) dpageUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		_, vy := v.Size()
		if oy == 0 || oy+cy < vy {
			if err := v.SetOrigin(ox, 0); err != nil {
				return err
			}
		} else if oy <= vy {
			if err := v.SetOrigin(ox, oy+cy-vy); err != nil {
				return err
			}
		} else if err := v.SetOrigin(ox, oy-vy); err != nil {
			return err
		}
		if err := v.SetCursor(cx, 0); err != nil {
			return err
		}
	}
	return nil
}

func (gui *Gui) initStashedView(r *git.Repository) error {
	v, err := gui.g.View(stashedViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	_, cy := v.Cursor()
	_, oy := v.Origin()
	stashedItems := r.Stasheds
	for i, stashedItem := range stashedItems {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(v, "%s%d %s: %s (%s)\n", prefix, stashedItem.StashID, cyan.Sprint(stashedItem.BranchName), stashedItem.Description, cyan.Sprint(stashedItem.Hash))
	}
	return nil
}

// there is no AI, only too much if clauses
func (gui *Gui) initFocusStat(r *git.Repository) error {
	v, err := gui.g.View(detailViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	v.Title = " Status "
	fmt.Fprintln(v, "On branch "+cyan.Sprint(r.State.Branch.Name))
	ps, err := strconv.Atoi(r.State.Branch.Pushables)
	pl, er2 := strconv.Atoi(r.State.Branch.Pullables)
	// TODO: move to text-render
	if err != nil || er2 != nil {
		fmt.Fprintln(v, "Your branch is not tracking a remote branch.")
	} else {
		if ps == 0 && pl == 0 {
			fmt.Fprintln(v, "Your branch is up to date with "+cyan.Sprint(r.State.Remote.Branch.Name))
		} else {
			if ps > 0 && pl > 0 {
				fmt.Fprintln(v, "Your branch and "+cyan.Sprint(r.State.Remote.Branch.Name)+" have diverged,")
				fmt.Fprintln(v, "and have "+yellow.Sprint(r.State.Branch.Pushables)+" and "+yellow.Sprint(r.State.Branch.Pullables)+" different commits each, respectively.")
				fmt.Fprintln(v, "(\"pull\" to merge the remote branch into yours)")
			} else if pl > 0 && ps == 0 {
				fmt.Fprintln(v, "Your branch is behind "+cyan.Sprint(r.State.Remote.Branch.Name)+" by "+yellow.Sprint(r.State.Branch.Pullables)+" commit(s).")
				fmt.Fprintln(v, "(\"pull\" to update your local branch)")
			} else if ps > 0 && pl == 0 {
				fmt.Fprintln(v, "Your branch is ahead of "+cyan.Sprint(r.State.Remote.Branch.Name)+" by "+yellow.Sprint(r.State.Branch.Pushables)+" commit(s).")
				fmt.Fprintln(v, "(\"push\" to publish your local commits)")
			}
		}
	}
	files, err := command.Status(r)
	if err != nil {
		return err
	}
	stagedFiles = make([]*git.File, 0)
	unstagedFiles = make([]*git.File, 0)
	for _, file := range files {
		if file.X != git.StatusNotupdated && file.X != git.StatusUntracked && file.X != git.StatusIgnored && file.X != git.StatusUpdated {
			stagedFiles = append(stagedFiles, file)
		}
		if file.Y != git.StatusNotupdated {
			unstagedFiles = append(unstagedFiles, file)
		}
	}
	if len(stagedFiles) == 0 && len(unstagedFiles) == 0 {
		fmt.Fprintln(v, "\nNothing to commit, working tree clean")
	} else {
		if len(stagedFiles) > 0 {
			fmt.Fprintln(v, "\nChanges to be committed:")
			fmt.Fprintln(v, "")
			for _, f := range stagedFiles {
				fmt.Fprintln(v, " "+green.Sprint(string(f.X)+" "+f.Name))
			}
		}
		if len(unstagedFiles) > 0 {
			fmt.Fprintln(v, "\nChanges not staged for commit:")
			fmt.Fprintln(v, "")
			for _, f := range unstagedFiles {
				fmt.Fprintln(v, " "+red.Sprint(string(f.Y)+" "+f.Name))
			}
			fmt.Fprintln(v, "\n"+strconv.Itoa(len(stagedFiles))+" change(s) added to commit (use \"git add\")")
		}
		_, cy := v.Cursor()
		if cy == 0 {
			gui.focusStatusCursorDown(gui.g, v)
		}
	}
	return nil
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) focusStatusCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil && v.Title == " Status " {
		_, cy := v.Cursor()
		ox, oy := v.Origin()
		ly := v.BufferLines()

		ap := oy + cy
		// if we are at the end we just return
		if ap == len(ly)-2 {
			return nil
		}
		v.EditDelete(true)
		var next int
		for i := ap + 1; i < len(ly); i++ {
			if len(ly[i]) > 0 && ly[i][0] == ' ' {
				next = i - ap
				break
			}
		}
		if err := v.SetCursor(0, cy+next); err != nil {
			if err := v.SetOrigin(ox, oy+next); err != nil {
				return err
			}
		}
		v.EditWrite('→')
	}
	return nil
}

// moves the cursor upwards for the main view
func (gui *Gui) focusStatusCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil && v.Title == " Status " {
		ox, oy := v.Origin()
		_, cy := v.Cursor()
		ly := v.BufferLines()
		v.EditDelete(true)
		ap := oy + cy
		var prev int
		for i := ap - 1; i >= 0; i-- {
			if len(ly[i]) > 0 && ly[i][0] == ' ' {
				prev = ap - i
				break
			}
		}
		if err := v.SetCursor(0, cy-prev); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-prev); err != nil {
				return err
			}
		}
		v.EditWrite('→')
	}
	return nil
}

func (gui *Gui) addreset(g *gocui.Gui, v *gocui.View) error {
	if v.Title == " Status " {
		_, cy := v.Cursor()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		r := gui.getSelectedRepository()
		files, err := command.Status(r)
		if err != nil {
			return err
		}
		for _, f := range files {
			if strings.Contains(line, f.Name) {
				if f.X != git.StatusNotupdated && f.X != git.StatusUntracked && f.X != git.StatusIgnored && f.X != git.StatusUpdated {
					if err := command.Reset(r, f, &command.ResetOptions{}); err != nil {
						return err
					}
				}
				if f.Y != git.StatusNotupdated {
					if err := command.Add(r, f, &command.AddOptions{}); err != nil {
						return err
					}
					v.EditWrite('→')
				}

			}
		}
		return gui.initFocusStat(r)
	}
	return nil
}
