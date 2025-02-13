package store

import (
	"log"
	_ "net/http/pprof"
	"sort"
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
	data map[string]*URLNode
	head *URLNode
	tail *URLNode
}

type Store interface {
	update(url string, success bool, timeMs int64)
	filter(n int, sortBy string) []*URLNode
}

func New() {
	store := URLStore{
		data: make(map[string]*URLNode),
	}
	go processStoreRequests(store)
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
		if node.Prev != nil {
			node.Prev.Next = node.Next
		} else {
			s.head = node.Next
		}

		if node.Next != nil {
			node.Next.Prev = node.Prev
		} else {
			s.tail = node.Prev
		}

		if success {
			node.Data.Successes++
			node.Data.LastDownloadMs = timeMs
		} else {
			node.Data.Failures++
		}

		node.Data.LastSubmitted = time.Now()
		node.Data.Count++

		// Move node to end
		node.Prev, node.Next = s.tail, nil
		if s.tail != nil {
			s.tail.Next = node
		}
		s.tail = node

		return

	}

	// URL hasn't been submitted, request was successful, add it to the map
	if success {
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

		// Insert the new node at the tail
		if s.tail == nil {
			// if the list is empty, set head and tail
			s.head, s.tail = newNode, newNode
		} else {
			// List isn't empty, append the new node to the tail
			s.tail.Next = newNode
			newNode.Prev = s.tail
			s.tail = newNode
		}

		s.data[url] = newNode
	}

}

func (s *URLStore) filter(n int, sortBy string) []*URLNode {
	nodes := make([]*URLNode, 0, n)
	current := s.tail

	for current != nil && len(nodes) < n {
		nodes = append(nodes, current)
		current = current.Prev
	}

	// Sort by count, list is already sorted by newest to oldest
	if sortBy == "count" {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Data.Count > nodes[j].Data.Count
		})
	}

	return nodes
}
