package downloader

import (
	"spamhaus/store"
	"time"
)

type URL struct {
	URL  string
	Data store.URLData
}

var batchDownloader

func NewBatchWorkerPool() {
	workerPool := NewWorkerPool(3)
	batchDownloader = &Downloader{
		workerPool: workerPool,
		maxConcurrent: 3,
	}
}

func BatchProcess() {
	time.Sleep(60 * time.Second)
	topURLs := store.GlobalStore.GetTopURLs()
	for i := range topURLs {
		GlobalDownloader.DownloadURl(topURLs[i].URL)
	}
}

