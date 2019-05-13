package gui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/core/git"
	"github.com/isacikgoz/gitbatch/core/job"
	log "github.com/sirupsen/logrus"
)

var (
	black   = color.New(color.FgBlack)
	blue    = color.New(color.FgBlue)
	green   = color.New(color.FgGreen)
	red     = color.New(color.FgRed)
	cyan    = color.New(color.FgCyan)
	yellow  = color.New(color.FgYellow)
	white   = color.New(color.FgWhite)
	magenta = color.New(color.FgMagenta)

	bold = color.New(color.Bold)

	maxBranchLength     = 15
	maxRepositoryLength = 20
	hashLength          = 7

	ws       = " "
	pushable = string(blue.Sprint("↖"))
	pullable = string(blue.Sprint("↘"))
	dirty    = string(yellow.Sprint("✗"))

	queuedSymbol  = "•"
	workingSymbol = "•"
	successSymbol = "✔"
	pauseSymbol   = "॥"
	failSymbol    = "✗"

	fetchSymbol    = "↓"
	pullSymbol     = "↓↳"
	mergeSymbol    = "↳"
	checkoutSymbol = "↱"

	keySymbol = ws + yellow.Sprint("🔑") + ws

	sep = ws + yellow.Sprint("|") + ws

	modeSeperator       = ""
	keyBindingSeperator = "░"

	selectionIndicator = "→" + ws
	tab                = ws

	renderRules = &RepositoryDecorationRules{}
)

// RepositoryDecorationRules is a rule set for creating repositry labels
type RepositoryDecorationRules struct {
	MaxName        int
	MaxPushables   int
	MaxPullables   int
	MaxBranch      int
	MaxBranchDirty int
}

// this function handles the render and representation of the repository
// TODO: cleanup is required, right now it looks too complicated
func (gui *Gui) repositoryLabel(r *git.Repository) string {
	renderRules = gui.renderRules()

	gui.renderTableHeader(renderRules)

	var line string

	line = line + renderRevCount(r, renderRules) + sep
	line = line + renderBranchName(r, renderRules) + sep
	line = line + gui.renderRepoName(r, renderRules) + sep
	line = line + gui.renderStatus(r)

	return line
}

// render repo name, print green if cursor is on the repository
func (gui *Gui) renderRepoName(r *git.Repository, rule *RepositoryDecorationRules) string {
	var repoName string
	sr := gui.getSelectedRepository()
	if sr == r {
		repoName = selectionIndicator + green.Sprint(align(r.Name, rule.MaxName, true, true))
	} else {
		repoName = align(r.Name, rule.MaxName+2, true, true)
	}
	return repoName
}

// render branch, add x if it is dirty
func renderBranchName(r *git.Repository, rule *RepositoryDecorationRules) string {
	b := r.State.Branch
	branch := cyan.Sprint(b.Name)
	if !b.Clean {
		branch = align(branch+ws+dirty, rule.MaxBranchDirty, true, true)
	} else {
		branch = align(branch, rule.MaxBranch, true, true)
	}
	return branch
}

// render ahead and behind info
func renderRevCount(r *git.Repository, rule *RepositoryDecorationRules) string {
	var revCount string
	b := r.State.Branch
	if b.Pushables != "?" {
		revCount = pushable + ws + align(b.Pushables, rule.MaxPushables, false, false) +
			ws + pullable + ws + align(b.Pullables, rule.MaxPullables, false, false)
	} else {
		revCount = pushable + ws + yellow.Sprint(align(b.Pushables, rule.MaxPushables, false, false)) +
			ws + pullable + ws + yellow.Sprint(align(b.Pullables, rule.MaxPullables, false, false))
	}
	return revCount
}

// render working status of the repository
func (gui *Gui) renderStatus(r *git.Repository) string {
	var status string
	if r.WorkStatus() == git.Queued {
		if inQueue, j := gui.State.Queue.IsInTheQueue(r); inQueue {
			status = printQueued(r, j.JobType)
		}
	} else if r.WorkStatus() == git.Working {
		status = green.Sprint(workingSymbol) + ws + r.State.Message
	} else if r.WorkStatus() == git.Success {
		status = green.Sprint(successSymbol) + ws + r.State.Message
	} else if r.WorkStatus() == git.Paused {
		status = yellow.Sprint("! authentication required (u)")
	} else if r.WorkStatus() == git.Fail {
		status = red.Sprint(failSymbol) + ws + red.Sprint(r.State.Message)
	}
	return status
}

// render header of the table layout
func (gui *Gui) renderTableHeader(rule *RepositoryDecorationRules) {
	v, err := gui.g.View(mainViewFrameFeature.Name)
	if err != nil {
		log.Warn(err.Error())
	}
	v.Clear()
	var header string
	revlen := 2 + rule.MaxPullables + 2 + rule.MaxPushables + 1
	header = ws + magenta.Sprint(align("revs", revlen, true, true)) + sep
	header = header + align(magenta.Sprint("branch"), renderRules.MaxBranch, true, true) + sep
	header = header + magenta.Sprint(align("name", renderRules.MaxName+2, true, true)) + sep
	fmt.Fprintln(v, header)
}

// print queued item with the mode color
func printQueued(r *git.Repository, jt job.JobType) string {
	var info string
	switch jt {
	case job.FetchJob:
		info = blue.Sprint(queuedSymbol) + ws + "(" + blue.Sprint("fetch") + ws + r.State.Remote.Name + ")"
	case job.PullJob:
		info = magenta.Sprint(queuedSymbol) + ws + "(" + magenta.Sprint("pull") + ws + r.State.Remote.Name + ")"
	case job.MergeJob:
		info = cyan.Sprint(queuedSymbol) + ws + "(" + cyan.Sprint("merge") + ws + r.State.Branch.Upstream.Name + ")"
	default:
		info = green.Sprint(queuedSymbol)
	}
	return info
}

// render commit lable according to its status(local/even/remote)
func commitLabel(c *git.Commit, sel bool) string {
	re := regexp.MustCompile(`\r?\n`)
	msg := re.ReplaceAllString(c.Message, " ")
	if sel {
		msg = green.Sprint(msg)
	}
	var body string
	switch c.CommitType {
	case git.EvenCommit:
		body = cyan.Sprint(c.Hash[:hashLength]) + " " + msg
	case git.LocalCommit:
		body = blue.Sprint(c.Hash[:hashLength]) + " " + msg
	case git.RemoteCommit:
		if len(c.Hash) > hashLength {
			body = yellow.Sprint(c.Hash[:hashLength]) + " " + msg
		} else {
			body = yellow.Sprint(c.Hash[:len(c.Hash)]) + " " + msg
		}
	default:
		body = c.Hash[:hashLength] + " " + msg
	}
	return body
}

// limit the text length for visual concerns
func adjustTextLength(text string, maxLength int) string {
	if len(text) > maxLength {
		return text[:maxLength-2] + ".."
	}
	return text
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

// the remote link can be too verbose sometimes, so it is good to trim it
func trimRemoteURL(url string) (urltype string, shorturl string) {
	// lets trim the unnecessary .git extension of the url
	regit := regexp.MustCompile(`.git`)
	if regit.MatchString(url[len(url)-4:]) {
		url = url[:len(url)-4]
	}

	// find out the protocol
	ressh := regexp.MustCompile(`git@`)
	rehttp := regexp.MustCompile(`http://`)
	rehttps := regexp.MustCompile(`https://`)

	// separate the protocol and remote link
	if ressh.MatchString(url) {
		shorturl = ressh.Split(url, 5)[1]
		urltype = "ssh"
	} else if rehttp.MatchString(url) {
		shorturl = rehttp.Split(url, 5)[1]
		urltype = "http"
	} else if rehttps.MatchString(url) {
		shorturl = rehttps.Split(url, 5)[1]
		urltype = "https"
	}
	return urltype, shorturl
}

// DiffStatDecorationRules is a rule set for creating diffstat text
type DiffStatDecorationRules struct {
	MaxNameLength        int
	MaxChangeCountLength int
	MaxChangesLength     int
}

// DiffStatItem is a line of a diff stat
type DiffStatItem struct {
	FileName    string
	ChangeCount string
	Changes     string
}

// get output of "git show <commit> --shortstat" and convert it to DiffStatItem
// slice and generate rules
func genDiffStat(in string) (*DiffStatDecorationRules, []*DiffStatItem) {
	rules := &DiffStatDecorationRules{}
	stats := make([]*DiffStatItem, 0)

	re := regexp.MustCompile(`\s+\|\s+`)
	r1 := regexp.MustCompile(`\d+\s+`)

	for _, line := range strings.Split(in, "\n") {
		s := re.Split(line, 2)
		ds := &DiffStatItem{}
		ds.FileName = s[0]

		if rules.MaxNameLength < len(ds.FileName) {
			rules.MaxNameLength = len(ds.FileName)
		}

		if len(s) > 1 && r1.MatchString(s[1]) {
			cc := r1.FindString(s[1])
			ds.ChangeCount = strings.TrimSpace(cc)
			if rules.MaxChangeCountLength < len(ds.ChangeCount) {
				rules.MaxChangeCountLength = len(ds.ChangeCount)
			}
			d := r1.Split(s[1], 2)

			ds.Changes = d[1]
			if rules.MaxChangesLength < len(ds.Changes) {
				rules.MaxChangesLength = len(ds.Changes)
			}
		}
		stats = append(stats, ds)
	}
	return rules, stats
}

// colorize diff stat
func decorateDiffStat(in string, sum bool) string {
	var d string

	s := strings.Split(in, "\n")
	if sum {
		d = strconv.Itoa(len(s)-1) + " file(s) changed." + "\n\n"
	}
	rule, stats := genDiffStat(in)
	for _, stat := range stats {
		if len(stat.FileName) <= 0 {
			continue
		}
		d = d + cyan.Sprint(align(stat.FileName, rule.MaxNameLength, true, true)) + yellow.Sprint(" | ") + align(stat.ChangeCount, rule.MaxChangeCountLength, false, false) + " "
		sr := []rune(stat.Changes)
		for _, r := range sr {
			if r == '+' {
				d = d + green.Sprint(string(r))
			} else if r == '-' {
				d = d + red.Sprint(string(r))
			} else {
				d = d + string(r)
			}
		}
		d = d + "\n"
	}
	return d
}

// align text with whitespaces
func align(in string, max int, direction, trim bool) string {
	realmax := 50
	il := len(in)
	if max > realmax {
		max = 50
	}
	if trim && il > realmax {
		return " ..." + in[il-46:]
	}
	if il < max {
		if direction {
			in = in + strings.Repeat(" ", max-il)
		} else {
			in = strings.Repeat(" ", max-il) + in
		}
	}
	return in
}

// colorize commit info
func decorateCommit(in string) string {
	var d string
	lines := strings.Split(in, "\n")
	d = d + strings.Replace(lines[0], "Hash:", cyan.Sprint("Hash:"), 1) + "\n"
	d = d + strings.Replace(lines[1], "Author:", cyan.Sprint("Author:"), 1) + "\n"
	d = d + strings.Replace(lines[2], "Date:", cyan.Sprint("Date:"), 1) + "\n"
	for _, l := range lines[3:] {
		d = d + l + "\n"
	}
	d = d[:len(d)]
	return d
}

// repository render rules
func (gui *Gui) renderRules() *RepositoryDecorationRules {
	rules := &RepositoryDecorationRules{}

	for _, r := range gui.State.Repositories {
		if len(r.State.Branch.Pullables) > rules.MaxPullables {
			rules.MaxPullables = len(r.State.Branch.Pullables)
		}
		if len(r.State.Branch.Pushables) > rules.MaxPushables {
			rules.MaxPushables = len(r.State.Branch.Pushables)
		}
		if len(r.State.Branch.Name) > rules.MaxBranch {
			rules.MaxBranch = len(r.State.Branch.Name)
		}
		if len(r.Name) > rules.MaxName {
			rules.MaxName = len(r.Name)
		}
	}
	rules.MaxBranchDirty = rules.MaxBranch + len(cyan.Sprint("")) + len(dirty) + len(ws)
	rules.MaxBranch = rules.MaxBranch + len(cyan.Sprint("")) + 2
	return rules
}
