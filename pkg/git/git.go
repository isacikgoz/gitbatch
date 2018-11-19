package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"regexp"
)

type RepoEntity struct {
	Name       string
	AbsPath    string
	Repository git.Repository
	Pushables  string
	Pullables  string
	Branch     string
}

func InitializeRepository(directory string) (RepoEntity, error) {
	var entity RepoEntity
	file, err := os.Open(directory)
	if err != nil {
		return entity, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return entity, err
	}
	r, err := git.PlainOpen(directory)
	if err != nil {
		return entity, err
	}
	pushable, pullable := UpstreamDifferenceCount(directory)
	branch, err := CurrentBranchName(directory)
	entity = RepoEntity{fileInfo.Name(), directory, *r, pushable, pullable, branch}
	
	return entity, nil
}

func (entity *RepoEntity) GetCommits() (commits []string, err error) {
	r := entity.Repository
	//TODO: Handle Errors
	ref, _ := r.Head()

    cIter, _ := r.Log(&git.LogOptions{
    	From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})

// ... just iterates over the commits
    err = cIter.ForEach(func(c *object.Commit) error {
    	commitstring := string([]rune(c.Hash.String())[:7]) + " " + c.Message
    	re := regexp.MustCompile(`\r?\n`)
    	commitstring = re.ReplaceAllString(commitstring, " ")
        commits = append(commits, commitstring)

        return nil
	})
	if err != nil {
		return commits, err
	}
    return commits, nil
}

func (entity *RepoEntity) GetRemotes() (remotes []string, err error) {
	r := entity.Repository

    if list, err := r.Remotes(); err != nil {
        return remotes, err
    } else {
        for _, r := range list {
        	remoteString := r.Config().Name + " → " + r.Config().URLs[0]
            remotes = append(remotes, remoteString)
        }
    }
    return remotes, nil
}

func (entity *RepoEntity) GetStatus() (status string) {
	status = "↑ " + entity.Pushables + " ↓ " + entity.Pullables + " → " + entity.Branch
	re := regexp.MustCompile(`\r?\n`)
    status = re.ReplaceAllString(status, " ")
    return status
}



