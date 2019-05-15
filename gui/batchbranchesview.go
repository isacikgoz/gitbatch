package gui

import (
	"fmt"
	"sort"

	"github.com/jroimartin/gocui"
)

// close confirmation view
func (gui *Gui) openBatchBranchView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop(batchBranchViewFeature.Name); err != nil {
		return err
	}
	gui.renderBatchBranches()
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
func (gui *Gui) renderBatchBranches() error {
	v, err := gui.g.View(batchBranchViewFeature.Name)
	if err != nil {
		return err
	}
	v.Clear()
	if len(gui.State.Repositories) == 0 {
		return nil
	}
	branchMap := make(map[string]int, 0)
	for _, r := range gui.State.Repositories {
		for _, b := range r.Branches {
			branchMap[b.Name] = branchMap[b.Name] + 1
		}
	}
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range branchMap {
		ss = append(ss, kv{k, v})
	}
	var vals []int
	for _, val := range branchMap {
		vals = append(vals, val)
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	si := 0
	if len(gui.State.targetBranch) == 0 {
		gui.State.targetBranch = ss[0].Key
	}
	for i, kv := range ss {
		if kv.Key == gui.State.targetBranch {
			si = i
			fmt.Fprintf(v, "%s%s %d\n", ws, green.Sprint(kv.Key), kv.Value)
		} else {
			fmt.Fprintf(v, "%s%s %d\n", tab, kv.Key, kv.Value)
		}
		gui.State.totalBranches = append(gui.State.totalBranches, kv.Key)
	}
	adjustAnchor(si, len(ss), v)
	return nil
}
