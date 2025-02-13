package downloader

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"spamhaus/store"
	"sync"
	"time"
)

type WorkerPool struct {
	wg sync.WaitGroup
}

var (
	Requests = make(chan string)
	finished = make(chan struct{})
)

func NewWorkerPool(poolSize int) *WorkerPool {
	pool := &WorkerPool{}

	for i := 0; i < poolSize; i++ {
		go pool.worker()
	}

	return pool
}

// Shutdown closes the Requests channel to prevent more requests coming in
// then blocks on the finished channel
func (wp *WorkerPool) Shutdown() {
	log.Println("workerpool: attempting graceful shutdown")
	close(Requests)
	<-finished
	log.Println("workerpool: shutdown complete")
}

func AddTask(url string) {
	log.Printf("adding download task to worker pool URL: %s", url)
	Requests <- url
}

func (wp *WorkerPool) worker() {
	for url := range Requests {
		wp.wg.Add(1)
		start := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("worker pool error: downloading %s, %v", url, err)
			return
		}
		duration := time.Since(start).Milliseconds()
		defer resp.Body.Close()

		_, err = io.Copy(ioutil.Discard, resp.Body)
		store.Update(url, err == nil && resp.StatusCode == 200, duration)

		wp.wg.Done()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
