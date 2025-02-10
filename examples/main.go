package main

import (
	"fmt"
	"log"
	"time"

	"github.com/frhan23/jobber-go/job"
)

// ExampleJob represents a simple job structure
type ExampleJob struct {
	name string
}

// Execute runs the job logic
func (e ExampleJob) Execute() {
	log.Printf("[Job] Executing: %s\n", e.name)
	time.Sleep(1 * time.Second) // Simulate processing time
	log.Printf("[Job] Completed: %s\n", e.name)
}

func main() {
	fmt.Println("=== Running Sequential Job Queue ===")
	sequentialQueue := job.NewJobQueue(5) // Buffer size = 5

	// Enqueue some jobs sequentially
	sequentialQueue.Enqueue(ExampleJob{name: "Job 1"})
	sequentialQueue.Enqueue(ExampleJob{name: "Job 2"})
	sequentialQueue.Enqueue(ExampleJob{name: "Job 3"})
	sequentialQueue.Enqueue(ExampleJob{name: "Job 4"})
	sequentialQueue.Enqueue(ExampleJob{name: "Job 5"})

	// Close the queue after jobs are enqueued
	sequentialQueue.Close()
	fmt.Println("=== Sequential Queue Finished ===")

	// Pause before running async queue
	time.Sleep(2 * time.Second)

	fmt.Println("\n=== Running Asynchronous Job Queue ===")
	asyncQueue := job.NewAsyncJobQueue(5, 3) // Buffer = 5, workers =3

	// Enqueue some jobs asynchronously
	for i := 1; i <= 5; i++ {
		asyncQueue.Enqueue(ExampleJob{name: fmt.Sprintf("Async Job %d", i)})
	}

	// Close the queue gracefully
	time.Sleep(3 * time.Second) // Allow some jobs to process
	asyncQueue.Close()
	fmt.Println("=== Async Queue Finished ===")
}
