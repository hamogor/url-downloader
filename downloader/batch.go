package downloader

import (
	"log"
	"spamhaus/store"
	"time"
)

type BatchProcess struct {
	workerPool   *WorkerPool
	concurrency  int
	interval     time.Duration
	numberOfURLs int
}

var batchProcessor BatchProcess

func NewBatchProcess(interval time.Duration, poolSize, numberOfURLS int) {
	workerPool := NewWorkerPool(3)
	batchProcessor = BatchProcess{
		workerPool:   workerPool,
		concurrency:  poolSize,
		interval:     interval,
		numberOfURLs: numberOfURLS,
	}
	batchProcessor.Run()
}

func (b *BatchProcess) Run() {
	go func() {
		for {
			b.runJob()
			time.Sleep(time.Second * b.interval)
		}
	}()
}

func (b *BatchProcess) runJob() {
	log.Println("batch: starting batch process")

	topURLs := store.Filter(b.numberOfURLs, "")
	if len(topURLs) == 0 {
		log.Println("batch: no urls to process")
		return
	}

	for _, url := range topURLs {
		AddTask(url.URL)
	}

	b.workerPool.Wait()
	log.Println("batch: finished batch process")
	b.logStats(topURLs)
}

func (b *BatchProcess) logStats(topURLS []*store.URLNode) {
	log.Println("----- Batch Job Stats -----")

	if len(topURLS) == 0 {
		log.Println("batch: no urls processed in this batch")
		return
	}

	for _, node := range topURLS {
		data := node.Data
		log.Printf("URL: %s | Count: %d | Successes: %d | Failures: %d | Last Download Time: %dms",
			node.URL, data.Count, data.Successes, data.Failures, data.LastDownloadMs)
	}

	log.Println("----------------")
}
