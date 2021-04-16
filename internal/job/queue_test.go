package job

import (
	"testing"

	"github.com/isacikgoz/gitbatch/internal/git"
	"github.com/stretchr/testify/require"
)

func TestCreateJobQueue(t *testing.T) {
	if output := CreateJobQueue(); output == nil {
		t.Errorf("Test Failed.")
	}
}

func TestAddJob(t *testing.T) {
	th := git.InitTestRepositoryFromLocal(t)
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
	th := git.InitTestRepositoryFromLocal(t)
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
	th := git.InitTestRepositoryFromLocal(t)
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
	th := git.InitTestRepositoryFromLocal(t)
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
