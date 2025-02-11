package store

import (
	"log"
	"sort"
	"sync"
	"time"
)

var GlobalStore *Store

type URLNode struct {
	URL  string
	Data *URLData
	Prev *URLNode
	Next *URLNode
}

type Store struct {
	mu   *sync.RWMutex
	data map[string]*URLNode
	head *URLNode
	tail *URLNode
}

type URLStore interface {
	GetLatestURLs(n int, sortBy string) []map[string]*URLData
	UpdateURL(url string, success bool, timeMs int64)
	GetStats() map[string]*URLData
}

func New() {
	GlobalStore = &Store{
		mu:   &sync.RWMutex{},
		data: make(map[string]*URLNode),
	}
}

func (s *Store) GetLatestURLs(n int, sortBy string) []*URLNode {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]*URLNode, 0)
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

func (s *Store) GetTopURLs(n int) []*URLNode {
	s.mu.Lock()
	defer s.mu.Unlock()
	var nodes []*URLNode
	for node := s.head; node != nil; node = node.Next {
		nodes = append(nodes, node)
	}

	// Sort the slice of nodes by Count in descending order without changing the linked list
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Data.Count > nodes[j].Data.Count
	})

	// Return  the top 10 URLs
	if len(nodes) > n {
		nodes = nodes[:n]
	}

	return nodes
}

func (s *Store) UpdateURL(url string, success bool, timeMs int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *Store) GetStats() map[string]*URLNode {
	return s.data
}
