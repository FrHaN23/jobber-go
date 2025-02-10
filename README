# Queue-Go

A lightweight and efficient job queue system in Go, supporting both **sequential** and **asynchronous** processing using goroutines.

## 🚀 Features
- **Sequential Processing**: Jobs are executed one after another.
- **Asynchronous Processing**: Jobs are processed concurrently using multiple workers.
- **Graceful Shutdown**: Ensures all jobs finish before closing the queue.
- **Thread-Safe**: Uses mutexes to prevent race conditions.

## 📦 Installation
```sh
go get github.com/frhan23/jobber-go
```

## 🔧 Usage

### 1️⃣ **Sequential Job Queue**
```go
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
	queue := job.NewJobQueue(5)
	queue.Enqueue(ExampleJob{name: "Job 1"})
	queue.Enqueue(ExampleJob{name: "Job 2"})
	queue.Close()
}
```

### 2️⃣ **Asynchronous Job Queue**
```go
func main() {
	asyncQueue := job.NewJobQueue(5) // Buffer = 5
	
	for i := 1; i <= 5; i++ {
		asyncQueue.EnqueueAsync(ExampleJob{name: fmt.Sprintf("Async Job %d", i)})
	}
	
	time.Sleep(3 * time.Second)
	asyncQueue.Close()
}
```

## ✅ Running Tests
To run tests with verbose output and benchmarking:
```sh
go test -v -bench . ./job
```

To check code coverage:
```sh
go test -cover ./job
```

## 📜 License
This project is licensed under the MIT License.

