package downloader

import (
	"spamhaus/store"
	"sync"
)

type WorkerPool struct {
	wg    sync.WaitGroup
	queue chan string
}

func NewWorkerPool(poolSize int) *WorkerPool {
	pool := &WorkerPool{
		queue: make(chan string, poolSize),
	}

	for i := 0; i < poolSize; i++ {
		go pool.worker()
	}

	return pool
}

func (wp *WorkerPool) addTask(url string) {
	wp.wg.Add(1)
	wp.queue <- url
}

func (wp *WorkerPool) worker() {
	for url := range wp.queue {
		success, duration := fetchURL(url)
		if success {
			store.GlobalStore.UpdateURL(url, success, duration)
		}
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
