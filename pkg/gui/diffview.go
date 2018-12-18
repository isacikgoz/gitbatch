package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

var diffReturnView string

// renders the diff view
func (gui *Gui) prepareDiffView(g *gocui.Gui, v *gocui.View, display []string) (out *gocui.View, err error) {
	maxX, maxY := g.Size()
	diffReturnView = v.Name()
	out, err = g.SetView(diffViewFeature.Name, 5, 3, maxX-5, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return out, err
		}
	}
	out.Title = diffViewFeature.Title
	out.Overwrite = true
	out.Wrap = true
	if err = gui.focusToView(diffViewFeature.Name); err != nil {
		return out, err
	}
	for _, line := range display {
		fmt.Fprintln(out, line)
	}
	return out, err
}

// open diff view for the selcted commit
// called from commitview, so initial view is commitview
func (gui *Gui) openCommitDiffView(g *gocui.Gui, v *gocui.View) (err error) {
	entity := gui.getSelectedRepository()
	commit := entity.Commit
	commitDetail := []string{("Hash: " + cyan.Sprint(commit.Hash) + "\n" + "Author: " + commit.Author +
		"\n" + commit.Time + "\n" + "\n" + "\t\t" + commit.Message + "\n")}
	diff, err := git.Diff(entity, entity.Commit.Hash)
	if err != nil {
		return err
	}
	colorized := colorizeDiff(diff)
	commitDetail = append(commitDetail, colorized...)
	out, err := gui.prepareDiffView(g, v, commitDetail)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	out.Title = " Commit Detail "
	return nil
}

// called from status, so initial view may be stagedview or unstaged view
func (gui *Gui) openFileDiffView(g *gocui.Gui, v *gocui.View) (err error) {
	entity := gui.getSelectedRepository()
	_, cy := v.Cursor()
	_, oy := v.Origin()
	var files []*git.File
	switch v.Name() {
	case unstageViewFeature.Name:
		_, files, err = populateFileLists(entity)
	case stageViewFeature.Name:
		files, _, err = populateFileLists(entity)
	}
	if err != nil {
		return err
	}
	if len(files) <= 0 {
		return nil
	}
	output, err := files[cy+oy].Diff()
	if err != nil || len(output) <= 0 {
		return nil
	}
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			diffReturnView); err != nil {
			return err
		}
	}
	colorized := colorizeDiff(output)
	_, err = gui.prepareDiffView(g, v, colorized)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

// called from stashview, so initial view is stashview
func (gui *Gui) showStash(g *gocui.Gui, v *gocui.View) (err error) {
	entity := gui.getSelectedRepository()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if len(entity.Stasheds) <= 0 {
		return nil
	}
	stashedItem := entity.Stasheds[oy+cy]
	output, err := stashedItem.Show()
	if err != nil {
		if err = gui.openErrorView(g, output,
			"You should manually resolve this issue",
			stashViewFeature.Name); err != nil {
			return err
		}
	}
	colorized := colorizeDiff(output)
	_, err = gui.prepareDiffView(g, v, colorized)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

// close the opened diff view
func (gui *Gui) closeCommitDiffView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return nil
	}
	return gui.closeViewCleanup(diffReturnView)
}
