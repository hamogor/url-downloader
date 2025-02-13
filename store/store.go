package store

import (
	"container/heap"
	"log"
	_ "net/http/pprof"
	"time"
)

var (
	Requests = make(chan Request, 100)
	finished = make(chan struct{})
)

type Request struct {
	Method   string
	URL      string
	SortBy   string
	Response chan Response
	TimeMs   int64
	Number   int
	Success  bool
}

type Response struct {
	Output interface{}
}

type URLData struct {
	LastDownloadMs int64
	Count          int
	Successes      int
	Failures       int
	LastSubmitted  time.Time
}

type URLNode struct {
	URL  string
	Data *URLData
	Prev *URLNode
	Next *URLNode
}

type URLStore struct {
	data       map[string]*URLNode
	countHeap  *URLHeap
	latestHeap *URLHeap
}

type Store interface {
	update(url string, success bool, timeMs int64)
	filter(n int, sortBy string) []*URLNode
}

func New() {
	store := URLStore{
		data: make(map[string]*URLNode),
	}

	store.countHeap = &URLHeap{
		By: func(i, j *URLNode) bool {
			return i.Data.Count > j.Data.Count
		},
	}

	store.latestHeap = &URLHeap{
		By: func(i, j *URLNode) bool {
			return i.Data.LastSubmitted.After(j.Data.LastSubmitted)
		},
	}

	heap.Init(store.countHeap)
	heap.Init(store.latestHeap)
	store.reorderHeapsPeriodically()
	go processStoreRequests(store)
}

func (s *URLStore) reorderHeapsPeriodically() {
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Reorder heaps every 500ms
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Reorder both heaps
				heap.Init(s.latestHeap) // Rebuild the recency heap

				heap.Init(s.countHeap) // Rebuild the count heap
				log.Println("Both heaps reordered")
			}
		}
	}()
}

func Shutdown() {
	log.Println("store: attempting graceful shutdown")
	close(Requests)
	<-finished
	log.Println("store: shutdown complete")
}

func processStoreRequests(s URLStore) {
	for request := range Requests {
		switch request.Method {
		case "update":
			s.update(request.URL, request.Success, request.TimeMs)
			request.Response <- Response{Output: "ok"}
		case "filter":
			data := s.filter(request.Number, request.SortBy)
			request.Response <- Response{Output: data}
		}
	}
	close(finished)
}

func Update(url string, success bool, timeMs int64) interface{} {
	responseChan := make(chan Response)
	defer close(responseChan)

	Requests <- Request{
		Method:   "update",
		URL:      url,
		Success:  success,
		TimeMs:   timeMs,
		Response: responseChan,
	}

	response := <-responseChan
	return response.Output
}

func Filter(n int, sortBy string) []*URLNode {
	responseChan := make(chan Response)
	defer close(responseChan)

	Requests <- Request{
		Method:   "filter",
		Number:   n,
		SortBy:   sortBy,
		Response: responseChan,
	}

	response := <-responseChan
	return response.Output.([]*URLNode)
}

func (s *URLStore) update(url string, success bool, timeMs int64) {
	// If this URL has already been submitted, update the data
	if node, exists := s.data[url]; exists {
		log.Printf("updating existing url: %s", url)

		// Update the URL node data
		if success {
			node.Data.Successes++
			node.Data.LastDownloadMs = timeMs
		} else {
			node.Data.Failures++
		}

		node.Data.LastSubmitted = time.Now()
		node.Data.Count++

		// Re-heapify the heap after the update
		heap.Fix(s.countHeap, findIndex(s.countHeap, node))
		heap.Fix(s.latestHeap, findIndex(s.latestHeap, node))
	} else if success {
		// Add a new URL node if it is successful
		log.Printf("adding new url: %s", url)
		newNode := &URLNode{
			URL: url,
			Data: &URLData{
				Count:          1,
				Successes:      1,
				LastDownloadMs: timeMs,
				LastSubmitted:  time.Now(),
			},
		}

		// Insert the new node into the data map
		s.data[url] = newNode

		// Push the new node into both heaps
		heap.Push(s.countHeap, newNode)
		heap.Push(s.latestHeap, newNode)
	}
}

func (s *URLStore) filter(n int, sortBy string) []*URLNode {
	var heapToUse *URLHeap
	switch sortBy {
	case "count":
		heapToUse = s.countHeap
	case "latest":
		heapToUse = s.latestHeap
	}

	// Pre-allocate slice for top N results
	nodes := make([]*URLNode, 0, n)

	// Efficiently extract top N elements without calling heap.Pop
	for i := 0; i < n && heapToUse.Len() > 0; i++ {
		nodes = append(nodes, heapToUse.Nodes[0]) // Grab the root node
		heapToUse.Nodes = heapToUse.Nodes[1:]     // Remove it without heapify
	}

	return nodes
}
