package downloader

import (
	"log"
	"spamhaus/store"
	"sync"
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
			time.Sleep(b.interval)
		}
	}()
}

func (b *BatchProcess) runJob() {
	log.Printf("batch: starting batch job")

	topURLS := store.GlobalStore.GetTopURLs(b.numberOfURLs)
	if len(topURLS) == 0 {
		log.Println("batch: no urls to process")
		return
	}

	var wg sync.WaitGroup
	queue := make(chan struct{}, b.concurrency)

	for _, node := range topURLS {
		time.Sleep(b.interval)
		wg.Add(1)
		queue <- struct{}{}

		go func(urlNode *store.URLNode) {
			defer wg.Done()
			defer func() { <-queue }()
			success, duration := fetchURL(urlNode.URL)
			store.GlobalStore.UpdateURL(urlNode.URL, success, duration)
		}(node)
	}

	wg.Wait()
	log.Println("batch: done. Download statistics:")
	b.logStats(topURLS)
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
