// Simple worker pool based on channels. Creating new goroutine if channel is full.
package workerpool

import (
	"sync"
)

type workerPool struct {
	tasks  chan Task
	closer *sync.Once
}

func New(workers int) workerPool {
	pool := workerPool{
		tasks:  make(chan Task, workers),
		closer: new(sync.Once),
	}
	for range workers {
		go pool.worker()
	}
	return pool
}

func (pool *workerPool) Stop() {
	pool.closer.Do(func() {
		close(pool.tasks)
	})
}

func (pool *workerPool) Go(task Task) {
	go func() { pool.tasks <- task }()
}

func (w *workerPool) worker() {
	for task := range w.tasks {
		if !task.IsCompleted() {
			task.Do()
		}
	}
}
