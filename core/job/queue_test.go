package job

import (
	"testing"

	"github.com/isacikgoz/gitbatch/core/git"
)

func TestCreateJobQueue(t *testing.T) {
	if output := CreateJobQueue(); output == nil {
		t.Errorf("Test Failed.")
	}
}

func TestAddJob(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	q := CreateJobQueue()
	var tests = []struct {
		input *Job
		same  bool
	}{
		{&Job{Repository: r}, false},
	}
	for _, test := range tests {
		if err := q.AddJob(test.input); err != nil && !test.same {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestRemoveFromQueue(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	q := CreateJobQueue()
	j := &Job{Repository: r}
	q.AddJob(j)
	var tests = []struct {
		input *git.Repository
	}{
		{r},
	}
	for _, test := range tests {
		if err := q.RemoveFromQueue(test.input); err != nil {
			t.Errorf("Test Failed. error: %s", err.Error())
		}
	}
}

func TestIsInTheQueue(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	q := CreateJobQueue()
	j := &Job{Repository: r}
	q.AddJob(j)
	var tests = []struct {
		input *git.Repository
	}{
		{r},
	}
	for _, test := range tests {
		if out1, out2 := q.IsInTheQueue(test.input); !out1 || j != out2 {
			t.Errorf("Test Failed. output: {%t, %s}", out1, out2.Repository.Name)
		}
	}
}

func TestStartJobsAsync(t *testing.T) {
	defer cleanRepo()
	r, err := testRepo()
	if err != nil {
		t.Fatalf("Test Failed. error: %s", err.Error())
	}
	q := CreateJobQueue()
	j := &Job{Repository: r}
	q.AddJob(j)
	var tests = []struct {
		input *JobQueue
	}{
		{q},
	}
	for _, test := range tests {
		if output := test.input.StartJobsAsync(); len(output) != 0 {
			t.Errorf("Test Failed.")
		}
	}
}
