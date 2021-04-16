package gui

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/jroimartin/gocui"
)

// close confirmation view
func (gui *Gui) openBatchBranchView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(batchBranchViewFeature.Name); err != nil {
		return err
	}
	_ = gui.renderBatchBranches(true)
	return gui.focusToView(batchBranchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) closeBatchBranchesView(g *gocui.Gui, v *gocui.View) error {
	if gui.order == focus {
		return nil
	}
	if _, err := g.SetViewOnBottom(batchBranchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(mainViewFeature.Name)
}

// updates the renderBatchBranches for given entity
func (gui *Gui) renderBatchBranches(calculate bool) error {
	v, err := gui.g.View(batchBranchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	if len(gui.State.Repositories) == 0 {
		return nil
	}
	branchMap := make(map[string]int)
	for _, r := range gui.State.Repositories {
		for _, b := range r.Branches {
			branchMap[b.Name] = branchMap[b.Name] + 1
		}
	}

	ss := make([]*branchCountMap, 0)
	for k, v := range branchMap {
		ss = append(ss, &branchCountMap{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Count == ss[j].Count {
			return lessAlphabetical(ss, i, j)
		}
		return ss[i].Count > ss[j].Count
	})
	if calculate {
		gui.State.totalBranches = ss
	}
	si := 0
	if len(gui.State.targetBranch) == 0 {
		gui.State.targetBranch = ss[0].BranchName
	}
	for i, kv := range gui.State.totalBranches {
		rule := gui.renderRules()
		n, branch := align(kv.BranchName, rule.MaxBranch, true)
		branch = branch + strings.Repeat(" ", n)
		if kv.BranchName == gui.State.targetBranch {
			si = i
			fmt.Fprintf(v, "%s%s%s%d\n", ws, green.Sprint(branch), sep, kv.Count)
		} else {
			fmt.Fprintf(v, "%s%s%s%d\n", tab, branch, sep, kv.Count)
		}
	}
	_ = adjustAnchor(si, len(gui.State.totalBranches), v)
	return nil
}

// Less is the interface implementation for Alphabetical sorting function
func lessAlphabetical(s []*branchCountMap, i, j int) bool {
	iRunes := []rune(s[i].BranchName)
	jRunes := []rune(s[j].BranchName)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}
	return false
}

// close confirmation view
func (gui *Gui) openSuggestBranchView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(suggestBranchViewFeature.Name); err != nil {
		return err
	}
	return gui.focusToView(suggestBranchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) closeSuggestBranchesView(g *gocui.Gui, v *gocui.View) error {
	if gui.order == focus {
		return nil
	}
	if _, err := g.SetViewOnBottom(suggestBranchViewFeature.Name); err != nil {
		return err
	}
	_ = gui.renderBatchBranches(false)
	return gui.focusToView(batchBranchViewFeature.Name)
}

// close confirmation view
func (gui *Gui) closeSuggestBranchesViewWithAdd(g *gocui.Gui, v *gocui.View) error {
	newBranch := strings.TrimSpace(v.ViewBuffer())
	if len(newBranch) <= 0 {
		return gui.closeSuggestBranchesView(g, v)
	}
	nm := &branchCountMap{BranchName: newBranch}
	gui.State.totalBranches = append(gui.State.totalBranches, nm)
	gui.State.targetBranch = nm.BranchName
	_ = gui.renderBatchBranches(false)
	return gui.closeSuggestBranchesView(g, v)
}
