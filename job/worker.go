package job

import (
	"log"
	"sync"
)

// Define the interface for jobs that can be submitted to queue
type Job interface {
	Execute()
}

// JobQueue represent a queue that process the job
type JobQueue struct {
	jobs  chan Job
	wg    sync.WaitGroup
	quit  chan struct{}
	close bool
	mutex sync.Mutex
}

// NewJobQueue Init the job
func NewJobQueue(bufferSize int) *JobQueue {
	jobQueue := &JobQueue{
		jobs:  make(chan Job, bufferSize),
		quit:  make(chan struct{}),
		wg:    sync.WaitGroup{},
		close: false,
		mutex: sync.Mutex{},
	}
	jobQueue.wg.Add(1)
	go jobQueue.worker()
	return jobQueue
}

func (jq *JobQueue) Enqueue(job Job) {
	jq.mutex.Lock()
	defer jq.mutex.Unlock()
	if jq.close {
		log.Print("[Queue] Cannot enqueue, queue is closed")
		return
	}
	select {
	case jq.jobs <- job:
		log.Print("[Queue] Job added to queue")
	default:
		log.Print("[Queue] Queue full, waiting for space")
		jq.jobs <- job // waiting for space
	}
}

func (jq *JobQueue) Close() {
	jq.mutex.Lock()
	if jq.close {
		jq.mutex.Unlock()
		return
	}
	jq.close = true
	close(jq.jobs)
	jq.mutex.Unlock()

	jq.wg.Wait()
	log.Print("Worker stopped")
}

func (jq *JobQueue) worker() {
	defer jq.wg.Done()
	log.Print("[Worker] Started")
	for job := range jq.jobs { // 🔥 No deadlock, exits when jobs is closed
		log.Println("[Worker] Processing job")
		job.Execute()
	}
	log.Print("[Worker] Stopped")
}
