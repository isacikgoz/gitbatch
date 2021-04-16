package git

import "fmt"

// Remote struct is simply a collection of remote branches and wraps it with the
// name of the remote and fetch/push urls. It also holds the *selected* remote
// branch
type Remote struct {
	Name     string
	URL      []string
	RefSpecs []string
	Branches []*RemoteBranch
}

// search for remotes in go-git way. It is the short way to get remotes but it
// does not give any insight about remote branches
func (r *Repository) initRemotes() error {
	rp := r.Repo
	r.Remotes = make([]*Remote, 0)

	rms, err := rp.Remotes()
	if err != nil {
		return err
	}
	for _, rm := range rms {
		rfs := make([]string, 0)
		for _, rf := range rm.Config().Fetch {
			rfs = append(rfs, string(rf))
		}
		remote := &Remote{
			Name:     rm.Config().Name,
			URL:      rm.Config().URLs,
			RefSpecs: rfs,
		}
		if err := remote.loadRemoteBranches(r); err != nil {
			continue
		}
		r.Remotes = append(r.Remotes, remote)
	}

	if len(r.Remotes) <= 0 {
		return fmt.Errorf("no remote for repository: %s", r.Name)
	}
	r.State.Remote = r.Remotes[0]
	return err
}
