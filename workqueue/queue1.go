package workqueue

import "sync"

type worker1 struct{}

func (worker *worker1) loopWork(jobCh <-chan JobRequest) {
	go func() {
		for {
			jobRequest := <-jobCh
			jobRequest.CurrentJob().Process()
			jobRequest.ResultChannel() <- nil
		}
	}()
}

// RunQueue1 creates a new work queue. This implementation has a simple job
// channel that delivers jobs to workers.
func RunQueue1(jobCh <-chan Job, workerCount int) (doneCh <-chan interface{}) {
	doneCh1 := make(chan interface{})
	jobRequestCh := make(chan JobRequest)
	workers := make([]*worker1, workerCount)

	for i := 0; i < workerCount; i++ {
		worker := &worker1{}
		worker.loopWork(jobRequestCh)
		workers[i] = worker
	}

	go func() {
		waitGroup := sync.WaitGroup{}

		for job := range jobCh {
			waitGroup.Add(1)

			go func(job Job) {
				resultCh := make(chan interface{})
				params := JobRequestParams{Job: job, ResultCh: resultCh}
				jobRequestCh <- NewJobRequest(params)
				<-resultCh
				waitGroup.Done()
			}(job)
		}

		waitGroup.Wait()
		doneCh1 <- nil
	}()

	doneCh = doneCh1
	return
}
