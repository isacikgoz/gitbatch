package job

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	mockJob1 := &Job{
		JobType:    PullJob,
		Repository: th.Repository,
	}
	mockJob2 := &Job{
		JobType:    FetchJob,
		Repository: th.Repository,
	}
	mockJob3 := &Job{
		JobType:    MergeJob,
		Repository: th.Repository,
	}

	var tests = []struct {
		input *Job
	}{
		{mockJob1},
		{mockJob2},
		{mockJob3},
	}
	for _, test := range tests {
		err := test.input.start()
		require.NoError(t, err)
	}
}
