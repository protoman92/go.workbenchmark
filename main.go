package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/protoman92/go.workbenchmark/workqueue"
)

func main() {
	jobCount := 1000
	workerCount := 10
	jobParams := workqueue.JobParams{ProcessDuration: time.Second}
	jobs := make([]workqueue.Job, jobCount)

	for i := 0; i < jobCount; i++ {
		jobs[i] = workqueue.NewJob(jobParams)
	}

	waitGroup := sync.WaitGroup{}

	runQueue := func(implementation int, runQueue func(<-chan workqueue.Job)) {
		waitGroup.Add(1)
		jobCh := make(chan workqueue.Job)

		go func() {
			for ix := range jobs {
				jobCh <- jobs[ix]
			}

			close(jobCh)
		}()

		start := time.Now()
		runQueue(jobCh)
		elapsed := time.Now().Sub(start)
		fmt.Printf("Elapsed time for queue %d: %v\n", implementation, elapsed)
		waitGroup.Done()
	}

	// Queue 1
	go func() {
		runQueue(1, func(jobCh <-chan workqueue.Job) {
			workqueue.RunQueue1(jobCh, workerCount)
		})
	}()

	// Queue 2
	go func() {
		runQueue(2, func(jobCh <-chan workqueue.Job) {
			workqueue.RunQueue2(jobCh, workerCount)
		})
	}()

	time.Sleep(time.Millisecond)
	waitGroup.Wait()
}
