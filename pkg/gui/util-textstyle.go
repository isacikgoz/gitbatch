package gui

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/gitbatch/pkg/git"
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

	fetchSymbol = "↓"
	pullSymbol  = "↓↳"
	mergeSymbol = "↳"

	keySymbol = ws + yellow.Sprint("🔑") + ws

	modeSeperator       = ""
	keyBindingSeperator = "░"

	selectionIndicator = ws + string(green.Sprint("→")) + ws
	tab                = ws
)

// this function handles the render and representation of the repository
// TODO: cleanup is required, right now it looks too complicated
func (gui *Gui) repositoryLabel(e *git.RepoEntity) string {

	var prefix string
	if e.Branch.Pushables != "?" {
		prefix = prefix + pushable + ws + e.Branch.Pushables +
			ws + pullable + ws + e.Branch.Pullables
	} else {
		prefix = prefix + pushable + ws + yellow.Sprint(e.Branch.Pushables) +
			ws + pullable + ws + yellow.Sprint(e.Branch.Pullables)
	}

	var repoName string
	se := gui.getSelectedRepository()
	if se == e {
		prefix = prefix + selectionIndicator
		repoName = green.Sprint(e.Name)
	} else {
		prefix = prefix + ws
		repoName = e.Name
	}
	// some branch names can be really long, in that times I hope the first
	// characters are important and meaningful
	branch := adjustTextLength(e.Branch.Name, maxBranchLength)
	prefix = prefix + string(cyan.Sprint(branch))

	if !e.Branch.Clean {
		prefix = prefix + ws + dirty + ws
	} else {
		prefix = prefix + ws
	}

	var suffix string
	// rendering the satus according to repository's state
	if e.State() == git.Queued {
		if inQueue, j := gui.State.Queue.IsInTheQueue(e); inQueue {
			switch mode := j.JobType; mode {
			case git.FetchJob:
				suffix = blue.Sprint(queuedSymbol)
			case git.PullJob:
				suffix = magenta.Sprint(queuedSymbol)
			case git.MergeJob:
				suffix = cyan.Sprint(queuedSymbol)
			default:
				suffix = green.Sprint(queuedSymbol)
			}
		}
		return prefix + repoName + ws + suffix
	} else if e.State() == git.Working {
		// TODO: maybe the type of the job can be written while its working?
		return prefix + repoName + ws + green.Sprint(workingSymbol)
	} else if e.State() == git.Success {
		return prefix + repoName + ws + green.Sprint(successSymbol)
	} else if e.State() == git.Paused {
		return prefix + repoName + ws + yellow.Sprint("authentication required (u)")
	} else if e.State() == git.Fail {
		return prefix + repoName + ws + red.Sprint(failSymbol) + ws + red.Sprint(e.Message)
	}
	return prefix + repoName
}

func commitLabel(c *git.Commit) string {
	var body string
	switch c.CommitType {
	case git.EvenCommit:
		body = cyan.Sprint(c.Hash[:hashLength]) + " " + c.Message
	case git.LocalCommit:
		body = blue.Sprint(c.Hash[:hashLength]) + " " + c.Message
	case git.RemoteCommit:
		body = yellow.Sprint(c.Hash[:hashLength]) + " " + c.Message
	default:
		body = c.Hash[:hashLength] + " " + c.Message
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
