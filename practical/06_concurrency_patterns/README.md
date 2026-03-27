# Worker Pool Pattern

The Worker Pool (or Thread Pool) pattern is essential for managing concurrency in Go, especially when:
1.  You have a large number of tasks to process.
2.  Creating a goroutine for every task would be too expensive (memory/CPU).
3.  You need to limit concurrency to avoid overwhelming downstream services (e.g., database limits).

## Key Components

### 1. Job Struct
Defines the unit of work. Often contains a `ResultCh` channel to send the result back to the submitter.
```go
type Job struct {
    ID       int
    ResultCh chan<- Result
}
```

### 2. Worker Pool Struct
Manages the lifecycle of workers.
```go
type WorkerPool struct {
    numWorkers int
    jobQueue   chan Job    // Buffered channel for pending jobs
    wg         sync.WaitGroup
    quit       chan struct{} // For graceful shutdown
}
```

### 3. Dispatcher / Start
Launches `numWorkers` goroutines. Each worker loops, consuming from `jobQueue`.

### 4. Graceful Shutdown
Crucial for robust systems. Two common strategies:
-   **Drain**: Close the `jobQueue`. Workers finish all pending jobs then exit.
-   **Hard Stop**: Close a `quit` channel. Workers finish *current* job then exit immediately (ignoring pending queue).

## Usage
```go
pool := NewWorkerPool(5, 100) // 5 workers, queue size 100
pool.Start()

pool.AddJob(myJob)
// ...
pool.StopAndDrain()
```
