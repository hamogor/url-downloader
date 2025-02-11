package downloader

import (
	"io"
	"log"
	"net/http"
	"time"
)

var GlobalDownloader *Downloader

type Downloader struct {
	workerPool *WorkerPool
}

func New(maxConcurrent int) {
	workerPool := NewWorkerPool(maxConcurrent)
	GlobalDownloader = &Downloader{
		workerPool: workerPool,
	}
}

func (d *Downloader) AddDownloadTask(url string) {
	log.Printf("starting download of %s", url)
	d.workerPool.addTask(url)
}

func fetchURL(url string) (bool, int64) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return false, 0
	}
	duration := time.Since(start).Milliseconds()
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body) // Discard response body
	if err == nil && resp.StatusCode == http.StatusOK {
		return true, duration
	}

	return false, 0
}
