package downloader

import "sync"

type WorkerPool struct {
	wg    sync.WaitGroup
	queue chan Task
}

func NewWorkerPool(poolSize int) *WorkerPool {
	pool := &WorkerPool{
		queue: make(chan Task, poolSize),
	}

	for i := 0; i < poolSize; i++ {
		go pool.worker()
	}

	return pool
}

func (wp *WorkerPool) addTask(task Task) {
	wp.wg.Add(1)
	wp.queue <- task
}

func (wp *WorkerPool) worker() {
	for task := range wp.queue {
		task.FetchURL()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
