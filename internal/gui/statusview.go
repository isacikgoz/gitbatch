package gui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/isacikgoz/gitbatch/internal/command"
	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/jroimartin/gocui"
)

// there is no AI, only too much if clauses
func (gui *Gui) initFocusStat(r *git.Repository) error {
	v, err := gui.g.View(dynamicViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	v.Title = string(StatusMode)
	if err := gui.updateDynamicKeybindings(); err != nil {
		return err
	}
	fmt.Fprintln(v, "On branch "+cyan.Sprint(r.State.Branch.Name))
	ps, err := strconv.Atoi(r.State.Branch.Pushables)
	pl, er2 := strconv.Atoi(r.State.Branch.Pullables)
	// TODO: move to text-render
	if err != nil || er2 != nil || r.State.Branch.Upstream == nil {
		fmt.Fprintln(v, "Your branch is not tracking a remote branch.")
	} else {
		if ps == 0 && pl == 0 {
			fmt.Fprintln(v, "Your branch is up to date with "+cyan.Sprint(r.State.Branch.Upstream.Name))
		} else {
			if ps > 0 && pl > 0 {
				fmt.Fprintln(v, "Your branch and "+cyan.Sprint(r.State.Branch.Upstream.Name)+" have diverged,")
				fmt.Fprintln(v, "and have "+yellow.Sprint(r.State.Branch.Pushables)+" and "+yellow.Sprint(r.State.Branch.Pullables)+" different commits each, respectively.")
				fmt.Fprintln(v, "(\"pull\" to merge the remote branch into yours)")
			} else if pl > 0 && ps == 0 {
				fmt.Fprintln(v, "Your branch is behind "+cyan.Sprint(r.State.Branch.Upstream.Name)+" by "+yellow.Sprint(r.State.Branch.Pullables)+" commit(s).")
				fmt.Fprintln(v, "(\"pull\" to update your local branch)")
			} else if ps > 0 && pl == 0 {
				fmt.Fprintln(v, "Your branch is ahead of "+cyan.Sprint(r.State.Branch.Upstream.Name)+" by "+yellow.Sprint(r.State.Branch.Pushables)+" commit(s).")
				fmt.Fprintln(v, "(\"push\" to publish your local commits)")
			}
		}
	}
	files, err := command.Status(r)
	if err != nil {
		return err
	}
	stagedFiles := make([]*git.File, 0)
	unstagedFiles := make([]*git.File, 0)
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
			fmt.Fprintln(v, "\n"+strconv.Itoa(len(stagedFiles))+" change(s) added to commit (consider \"add\")")
		}
		_, cy := v.Cursor()
		if cy == 0 {
			_ = gui.statusCursorDown(gui.g, v)
		} else {
			_ = gui.statusCursorUp(gui.g, v)
		}
	}
	return nil
}

// moves the cursor downwards for the main view and if it goes to bottom it
// prevents from going further
func (gui *Gui) statusCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
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
func (gui *Gui) statusCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, cy := v.Cursor()
		ly := v.BufferLines()
		v.EditDelete(true)
		ap := oy + cy
		var prev int
		for i := ap - 1; i >= 0; i-- {
			if i < len(ly) && len(ly[i]) > 0 && ly[i][0] == ' ' {
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

// add or reset file
func (gui *Gui) statusAddReset(g *gocui.Gui, v *gocui.View) error {
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
			}
		}
	}
	return gui.initFocusStat(r)
}

// show diff of the file
func (gui *Gui) statusDiff(g *gocui.Gui, v *gocui.View) error {

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
			out, err := command.DiffFile(f)
			if err != nil {
				v.Clear()
				v.Title = string(FileDiffMode)
				if err := gui.updateDynamicKeybindings(); err != nil {
					return err
				}
				fmt.Fprintln(v, "Can't get diff")
				return nil
			}
			v.Clear()
			v.Title = string(FileDiffMode)
			if err := gui.updateDynamicKeybindings(); err != nil {
				return err
			}
			fmt.Fprintln(v, strings.Join(colorizeDiff(out), "\n"))
		}
	}
	return nil
}

// stash uncommitted changes of the working directory
func (gui *Gui) stashChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	output, err := r.Stash()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	if err := gui.focusToRepository(g, v); err != nil {
		return err
	}
	if err := gui.initStashedView(r); err != nil {
		return err
	}
	return nil
}

// stats of current status, basically emulates "git status"
func (gui *Gui) statusStat(g *gocui.Gui, v *gocui.View) error {

	r := gui.getSelectedRepository()
	if err := gui.initFocusStat(r); err != nil {
		return err
	}
	return nil
}

// add all items to index
func (gui *Gui) statusAddAll(g *gocui.Gui, v *gocui.View) error {

	r := gui.getSelectedRepository()
	if err := command.AddAll(r, &command.AddOptions{}); err != nil {
		return err
	}
	if err := gui.initFocusStat(r); err != nil {
		return err
	}
	return nil
}

// reset all indexed items
func (gui *Gui) statusResetAll(g *gocui.Gui, v *gocui.View) error {

	r := gui.getSelectedRepository()
	ref, err := r.Repo.Head()
	if err != nil {
		return err
	}
	if err := command.ResetAll(r, &command.ResetOptions{
		Hash:        ref.Hash().String(),
		ResetType:   command.ResetMixed,
		CommandMode: command.ModeNative,
	}); err != nil {
		return err
	}
	if err := gui.initFocusStat(r); err != nil {
		return err
	}
	return nil
}
