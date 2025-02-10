package downloader

import (
	"log"
	"net/http"
	"spamhaus/store"
	"time"
)

var GlobalDownloader *Downloader

type Downloader struct {
	workerPool    *WorkerPool
	maxConcurrent int
}

type Task string

func New(maxConcurrent int, poolSize int) {
	workerPool := NewWorkerPool(poolSize)
	GlobalDownloader = &Downloader{
		workerPool:    workerPool,
		maxConcurrent: maxConcurrent,
	}
}

func (d *Downloader) DownloadURl(url string) {
	log.Printf("starting download of %s", url)
	d.workerPool.addTask(Task(url))
}

func (t Task) FetchURL() {
	startTime := time.Now()

	resp, err := http.Get(string(t))
	if err != nil {
		log.Printf("error: failed to fetch URL: %s, %v", t, err)
		// TODO update download stats to say it failed
		return
	}
	defer resp.Body.Close()

	duration := time.Since(startTime).Milliseconds()
	store.GlobalStore.UpdateURL(string(t), true, duration)

	data := store.GlobalStore.GetStats()
	log.Printf("%v", data)
	BatchProcess()

}
