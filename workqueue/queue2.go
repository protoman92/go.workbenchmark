package workqueue

import "sync"

type worker2 struct {
	jobQueue chan JobRequest
}

func (worker *worker2) loopWork(queue chan<- (chan JobRequest)) {
	go func() {
		for {
			queue <- worker.jobQueue

			select {
			case jobRequest, ok := <-worker.jobQueue:
				if !ok {
					return
				}

				jobRequest.CurrentJob().Process()
				jobRequest.ResultChannel() <- nil
			}
		}
	}()
}

// RunQueue2 creates a new work queue. This implementation has a channel of
// JobRequest channels that receive ready signals from workers, upon which it
// pushes jobs onto these channels for workers to process.
func RunQueue2(jobCh <-chan Job, workerCount int) {
	doneCh := make(chan interface{})
	jobQueueCh := make(chan chan JobRequest)
	workers := make([]*worker2, workerCount)

	for i := 0; i < workerCount; i++ {
		worker := &worker2{jobQueue: make(chan JobRequest)}
		worker.loopWork(jobQueueCh)
		workers[i] = worker
	}

	go func() {
		waitGroup := sync.WaitGroup{}

		for job := range jobCh {
			waitGroup.Add(1)
			jobQueue := <-jobQueueCh

			go func(job Job) {
				resultCh := make(chan interface{})
				params := JobRequestParams{Job: job, ResultCh: resultCh}
				jobQueue <- NewJobRequest(params)
				<-resultCh
				waitGroup.Done()
			}(job)
		}

		waitGroup.Wait()
		doneCh <- nil
	}()

	<-doneCh
}
