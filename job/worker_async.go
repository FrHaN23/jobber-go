package job

import (
	"log"
	"sync"
)

// AsyncJobQueue handles **concurrent job processing** using workers.
type AsyncJobQueue struct {
	jobs    chan Job
	wg      sync.WaitGroup
	quit    chan struct{}
	close   bool
	mutex   sync.Mutex
	workers int
}

// NewAsyncJobQueue initializes jobs queue with worker pool.
func NewAsyncJobQueue(bufferSize, workers int) *AsyncJobQueue {
	if bufferSize < workers { // Ensure the buffer is at least equal to worker count
		bufferSize = workers
	}
	asyncQueue := &AsyncJobQueue{
		jobs:    make(chan Job, bufferSize),
		quit:    make(chan struct{}),
		wg:      sync.WaitGroup{},
		close:   false,
		mutex:   sync.Mutex{},
		workers: workers,
	}

	// Spawn multiple workers
	for i := 0; i < workers; i++ {
		asyncQueue.wg.Add(1)
		go asyncQueue.worker(i)
	}

	return asyncQueue
}

// EnqueueAsync for **async processing using worker pool**.
func (jq *AsyncJobQueue) Enqueue(job Job) {
	jq.mutex.Lock()
	if jq.close {
		jq.mutex.Unlock()
		log.Print("[AsyncQueue] Cannot enqueue, queue is closed")
		return
	}
	jq.mutex.Unlock()

	select {
	case jq.jobs <- job:
		log.Print("[AsyncQueue] Job added")
	default:
		go func() { jq.jobs <- job }() // Offload to a goroutine instead of dropping
		log.Print("[AsyncQueue] Job queue full, waiting...")
	}
}

// Close gracefully shuts down queue.
func (jq *AsyncJobQueue) Close() {
	jq.mutex.Lock()
	if jq.close {
		jq.mutex.Unlock()
		return
	}
	jq.close = true
	close(jq.jobs)
	jq.mutex.Unlock()

	jq.wg.Wait()
	log.Print("[AsyncQueue] Workers stopped")
}

// worker for **async queue** (multiple workers).
func (jq *AsyncJobQueue) worker(id int) {
	defer jq.wg.Done()
	log.Printf("[Worker-%d] Started\n", id)

	for job := range jq.jobs {
		log.Printf("[Worker-%d] Processing job\n", id)
		job.Execute()
	}

	log.Printf("[Worker-%d] Stopped\n", id)
}
