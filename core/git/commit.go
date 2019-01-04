package git

import (
	"os/exec"
	"regexp"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Commit is the lightweight version of go-git's Reference struct. it holds
// hash of the commit, author's e-mail address, Message (subject and body
// combined) commit date and commit type wheter it is local commit or a remote
type Commit struct {
	Hash       string
	Author     string
	Message    string
	Time       string
	CommitType CommitType
}

// CommitType is the Type of the commit; it can be local or remote (upstream diff)
type CommitType string

const (
	// LocalCommit is the commit that not pushed to remote branch
	LocalCommit CommitType = "local"
	// EvenCommit is the commit that recorded locally
	EvenCommit CommitType = "even"
	// RemoteCommit is the commit that not merged to local branch
	RemoteCommit CommitType = "remote"
)

// NextCommit iterates over next commit of a branch
// TODO: the commits entites can tied to branch instead ot the repository
func (r *Repository) NextCommit() {
	r.State.Commit = r.Commits[(r.currentCommitIndex()+1)%len(r.Commits)]
}

// PreviousCommit iterates to opposite direction
func (r *Repository) PreviousCommit() {
	r.State.Commit = r.Commits[(len(r.Commits)+r.currentCommitIndex()-1)%len(r.Commits)]
}

// returns the active commit index
func (r *Repository) currentCommitIndex() int {
	cix := 0
	for i, c := range r.Commits {
		if c.Hash == r.State.Commit.Hash {
			cix = i
		}
	}
	return cix
}

// loads the local commits by simply using git log way. ALso, gets the upstream
// diff commits
func (r *Repository) loadCommits() error {
	rp := r.Repo
	r.Commits = make([]*Commit, 0)
	ref, err := rp.Head()
	if err != nil {
		log.Trace("Cannot get HEAD " + err.Error())
		return err
	}
	// git log first
	cIter, err := rp.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		log.Trace("git log failed " + err.Error())
		return err
	}
	defer cIter.Close()
	// find commits that fetched from upstream but not merged commits
	rmcs, _ := r.pullDiffsToUpstream()
	r.Commits = append(r.Commits, rmcs...)

	// find commits that not pushed to upstream
	lcs, _ := r.pushDiffsToUpstream()

	// ... just iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {
		re := regexp.MustCompile(`\r?\n`)
		cmType := EvenCommit
		for _, lc := range lcs {
			if lc == re.ReplaceAllString(c.Hash.String(), " ") {
				cmType = LocalCommit
			}
		}
		commit := &Commit{
			Hash:       re.ReplaceAllString(c.Hash.String(), " "),
			Author:     c.Author.Email,
			Message:    re.ReplaceAllString(c.Message, " "),
			Time:       c.Author.When.String(),
			CommitType: cmType,
		}
		r.Commits = append(r.Commits, commit)
		return nil
	})
	return err
}

// this function creates the commit entities according to active branchs diffs
// to *its* configured upstream
func (r *Repository) pullDiffsToUpstream() ([]*Commit, error) {
	remoteCommits := make([]*Commit, 0)
	pullables, err := RevList(r, RevListOptions{
		Ref1: "HEAD",
		Ref2: "@{u}",
	})
	if err != nil {
		// possibly found nothing or no upstream set
	} else {
		re := regexp.MustCompile(`\r?\n`)
		for _, s := range pullables {
			if len(s) < hashLength {
				continue
			}
			commit := &Commit{
				Hash:       s,
				Author:     gitShowEmail(r.AbsPath, s),
				Message:    re.ReplaceAllString(gitShowBody(r.AbsPath, s), " "),
				Time:       gitShowDate(r.AbsPath, s),
				CommitType: RemoteCommit,
			}
			remoteCommits = append(remoteCommits, commit)
		}
	}
	return remoteCommits, nil
}

// this function returns the hashes of the commits that are not pushed to the
// upstream of the specific branch
func (r *Repository) pushDiffsToUpstream() ([]string, error) {
	pushables, err := RevList(r, RevListOptions{
		Ref1: "@{u}",
		Ref2: "HEAD",
	})
	if err != nil {
		return make([]string, 0), nil
	}
	return pushables, nil
}

// gitShowEmail gets author's e-mail with git show command
func gitShowEmail(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ae", hash}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}

// gitShowBody gets body of the commit with git show
func gitShowBody(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%B", hash}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}

// gitShowDate gets commit's date with git show as string
func gitShowDate(repoPath, hash string) string {
	args := []string{"show", "--quiet", "--pretty=format:%ai", hash}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}
