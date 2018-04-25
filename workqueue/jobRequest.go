package workqueue

// JobRequestParams represents the required parameters to build a job request.
type JobRequestParams struct {
	Job      Job
	ResultCh chan interface{}
}

type jobRequest struct {
	JobRequestParams
}

// JobRequest represents a job request with a result channel.
type JobRequest interface {
	CurrentJob() Job
	ResultChannel() chan<- interface{}
}

// NewJobRequest creates a new job request.
func NewJobRequest(params JobRequestParams) JobRequest {
	return &jobRequest{JobRequestParams: params}
}

func (jobRequest *jobRequest) CurrentJob() Job {
	return jobRequest.Job
}

func (jobRequest *jobRequest) ResultChannel() chan<- interface{} {
	return jobRequest.ResultCh
}
