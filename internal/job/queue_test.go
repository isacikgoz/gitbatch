package job

import (
	"runtime"
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/isacikgoz/gitbatch/internal/testlib"
	"github.com/stretchr/testify/require"
)

func TestCreateJobQueue(t *testing.T) {
	if output := CreateJobQueue(); output == nil {
		t.Errorf("Test Failed.")
	}
}

func TestAddJob(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	q := CreateJobQueue()
	var tests = []struct {
		input *Job
	}{
		{&Job{Repository: th.Repository}},
	}
	for _, test := range tests {
		err := q.AddJob(test.input)
		require.NoError(t, err)
	}
}

func TestRemoveFromQueue(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	q := CreateJobQueue()
	j := &Job{Repository: th.Repository}
	err := q.AddJob(j)
	require.NoError(t, err)

	var tests = []struct {
		input *git.Repository
	}{
		{th.Repository},
	}
	for _, test := range tests {
		err := q.RemoveFromQueue(test.input)
		require.NoError(t, err)
	}
}

func TestIsInTheQueue(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	q := CreateJobQueue()
	j := &Job{Repository: th.Repository}
	err := q.AddJob(j)
	require.NoError(t, err)

	var tests = []struct {
		input *git.Repository
	}{
		{th.Repository},
	}
	for _, test := range tests {
		out1, out2 := q.IsInTheQueue(test.input)
		require.True(t, out1)
		require.Equal(t, j, out2)
	}
}

func TestStartJobsAsync(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	q := CreateJobQueue()
	j := &Job{Repository: th.Repository}
	err := q.AddJob(j)
	require.NoError(t, err)

	var tests = []struct {
		input *Queue
	}{
		{q},
	}
	for _, test := range tests {
		output := test.input.StartJobsAsync()
		require.Empty(t, output)
	}
}

func TestStartJobsAsyncManyJobs(t *testing.T) {
	th := testlib.InitTestRepositoryFromLocal(t)
	defer th.CleanUp(t)

	// enough jobs to keep every worker busy so they drain the queue concurrently
	n := 4 * runtime.GOMAXPROCS(0)
	q := CreateJobQueue()
	repos := make([]*git.Repository, 0, n)
	for i := 0; i < n; i++ {
		r, err := git.InitializeRepo(th.RepoPath)
		require.NoError(t, err)
		require.NoError(t, q.AddJob(&Job{Repository: r}))
		repos = append(repos, r)
	}

	require.Empty(t, q.StartJobsAsync())
	for _, r := range repos {
		inTheQueue, _ := q.IsInTheQueue(r)
		require.False(t, inTheQueue)
		require.Equal(t, git.Available, r.WorkStatus())
	}
}
