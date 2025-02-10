package job_test

import (
	"sync"
	"testing"
	"time"

	"github.com/frhan23/jobber-go/job"
)

type TestJob struct {
	id       int
	executed bool
	mutex    sync.Mutex
}

func (t *TestJob) Execute() {
	t.mutex.Lock()
	t.executed = true
	t.mutex.Unlock()
}

func TestJobQueueSequential(t *testing.T) {
	jq := job.NewJobQueue(2)
	defer jq.Close()
	job1 := &TestJob{id: 1}
	job2 := &TestJob{id: 2}
	job3 := &TestJob{id: 3}
	jq.Enqueue(job1)
	jq.Enqueue(job2)
	jq.Enqueue(job3)
	time.Sleep(100 * time.Millisecond)
	if !job1.executed || !job2.executed {
		t.Error("Expected job1 and job2 to be execute")
	}
	if !job3.executed {
		t.Error("Expected job3 to be execute")
	}
}

func TestJobQueueAsync(t *testing.T) {
	jq := job.NewJobQueue(2)
	defer jq.Close()

	job1 := &TestJob{id: 1}
	job2 := &TestJob{id: 2}
	job3 := &TestJob{id: 3}

	jq.EnqueueAsync(job1)
	jq.EnqueueAsync(job2)
	jq.EnqueueAsync(job3) // All should execute concurrently

	time.Sleep(50 * time.Millisecond) // Small delay to allow async jobs to finish

	if !job1.executed || !job2.executed || !job3.executed {
		t.Error("Expected job1, job2, and job3 to be executed asynchronously")
	}
}

func TestJobQueueClosure(t *testing.T) {
	jq := job.NewJobQueue(2)
	job1 := &TestJob{id: 1}
	jq.Enqueue(job1)

	// Wait for job execution before closing
	time.Sleep(100 * time.Millisecond)

	jq.Close()

	if !job1.executed {
		t.Error("Expected job1 to be executed")
	}

	// Attempt to enqueue after closing should return immediately
	done := make(chan bool)
	go func() {
		jq.Enqueue(&TestJob{id: 2})
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("Enqueue after close did not return promptly")
	}
}

func TestEnqueueAsyncAfterClose(t *testing.T) {
	jq := job.NewJobQueue(2) // Assume this is your async queue
	jq.Close()               // Close the queue first

	job := &TestJob{id: 1}
	jq.EnqueueAsync(job) // Try to enqueue after closing

	time.Sleep(50 * time.Millisecond) // Give time for log

	if job.executed {
		t.Error("Job should NOT be executed after queue is closed")
	}
}

func TestCloseTwice(t *testing.T) {
	jq := job.NewJobQueue(2)
	jq.Close()
	jq.Close() // should not panic
}

func BenchmarkJobQueue(b *testing.B) {
	jq := job.NewJobQueue(10)
	defer jq.Close()

	for i := 0; i < b.N; i++ {
		jq.EnqueueAsync(&TestJob{id: i})
	}
}
