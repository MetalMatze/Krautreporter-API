package workerqueue

//{
//	"processed": 0,
//	"failed": 0,
//	"jobs": {}
//}

// WorkerQueue holds all jobs that the workers wait for
type WorkerQueue struct {
	jobs      chan func()
	processed int
}

// NewWorkerQueue creates a WorkerQueue with n workers for your jobs
func New(n int) WorkerQueue {
	wq := WorkerQueue{
		jobs: make(chan func()),
	}

	// spawn n workers
	for i := 0; i < n; i++ {
		go wq.worker()
	}

	return wq
}

func (wq WorkerQueue) worker() {
	for j := range wq.jobs {
		j()
	}
}

// Push adds another job to the jobs chan to be processed by a worker
func (wq WorkerQueue) Push(j func()) {
	wq.processed = wq.processed + 1
	wq.jobs <- j
}
