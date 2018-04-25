package workqueue

import (
	"time"
)

// JobParams represents the required parameters to build a Job.
type JobParams struct {
	ProcessDuration time.Duration
}

type job struct {
	JobParams
}

// Job represents an executable job.
type Job interface {
	Process()
}

// NewJob creates a new Job.
func NewJob(params JobParams) Job {
	return &job{JobParams: params}
}

func (job *job) Process() {
	time.Sleep(job.ProcessDuration)
}
