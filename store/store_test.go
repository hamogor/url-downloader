package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func newStore(n int) {
	New()
	for i := 0; i < n; i++ {
		Update(
			fmt.Sprintf("http://example%d.com", i),
			true,
			int64(100+i),
		)
	}
}

// TestStore_GetLatestURLs ensures that the latest n urls are returned
// either by their time submitted or by their count
//func TestStore_GetLatestURLs(t *testing.T) {
//	newStore(15)
//
//	latest := Filter(5, "latest")
//
//	// This list is in the order we expect back from the store when getting by "latest"
//	expectedURLS := []string{
//		"http://example14.com",
//		"http://example13.com",
//		"http://example12.com",
//		"http://example11.com",
//		"http://example10.com",
//	}
//
//	for i := range latest {
//		if latest[i].URL != expectedURLS[i] {
//			t.Errorf("expected %s, got %s", expectedURLS[i], latest[i].URL)
//		}
//	}
//
//	// Update the first 5 URLs with some extra Counts
//	n := 1
//	for i := 0; i < 5; i++ {
//		for j := 0; j < n; j++ {
//			Update(
//				fmt.Sprintf("http://example1%d.com", i),
//				true,
//				int64(100+i),
//			)
//		}
//		n++
//	}
//
//	// The 5 urls we got should now have counts 6, 5, 4, 3, 2 in that order
//	count := Filter(5, "count")
//	for i, node := range count {
//		expectedCount := 6 - i
//		if node.Data.Count != expectedCount {
//			t.Errorf("expected %d, got %d", expectedCount, node.Data.Count)
//		}
//	}
//
//}

//// TestStore_GetTopURLs ensures that the top n counts on urls are returned
//func TestStore_GetTopURLs(t *testing.T) {
//	newStore(15)
//
//	// Update the first ten urls in the store with an extra counter
//	for i := 0; i < 10; i++ {
//		Update(
//			fmt.Sprintf("http://example%d.com", i),
//			true,
//			int64(100+i),
//		)
//	}
//
//	topURLS := Filter(10, "latest")
//
//	for i := 0; i < len(topURLS); i++ {
//		if topURLS[i].Data.Count != 2 {
//			t.Errorf("expected count of URL: %s to be 2, got %d", topURLS[i].URL, topURLS[i].Data.Count)
//		}
//	}
//}

// TestStore_UpdateURL tests that an added node is added to the end of the list
// and that an updated node is moved to the end of the list and count incremented
//func TestStore_UpdateURL(t *testing.T) {
//	tests := []struct {
//		name          string
//		url           string
//		success       bool
//		timeMs        int64
//		expectedHead  string
//		expectedTail  string
//		expectedCount int
//	}{
//		{
//			name:          "Add a new node",
//			url:           "http://example15.com",
//			success:       true,
//			timeMs:        100,
//			expectedHead:  "http://example0.com",
//			expectedTail:  "http://example15.com",
//			expectedCount: 1,
//		},
//		{
//			name:          "Update an existing node",
//			url:           "http://example1.com",
//			success:       true,
//			timeMs:        100,
//			expectedHead:  "http://example0.com",
//			expectedTail:  "http://example1.com",
//			expectedCount: 2,
//		},
//	}
//
//	store := URLStore{
//		mu:   &sync.RWMutex{},
//		data: make(map[string]*URLNode),
//	}
//	go processStoreRequests(store)
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			_ = Update(tt.url, tt.success, tt.timeMs)
//
//			// Check the head URL
//			if store.head.URL != tt.expectedHead {
//				t.Errorf("expected head %s, but got %s", tt.expectedHead, GlobalStore.head.URL)
//			}
//
//			// Check the tail URL
//			if store.tail.URL != tt.expectedTail {
//				t.Errorf("expected tail %s, but got %s", tt.expectedTail, GlobalStore.tail.URL)
//			}
//
//			// Check the count of the last node
//			if store.tail.Data.Count != tt.expectedCount {
//				t.Errorf("expected count %d, but got %d", tt.expectedCount, GlobalStore.tail.Data.Count)
//			}
//		})
//	}
//}

// Benchmark to test fetching the latest 50 URLs from the store
func BenchmarkGetLatestURLs(b *testing.B) {
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	log.SetOutput(ioutil.Discard)
	// Populate the store with 1000 URLs
	newStore(10000)

	b.ResetTimer()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Benchmark fetching the latest 50 URLs
	for i := 0; i < b.N; i++ {
		Filter(50, "")
	}
}

// Benchmark to test fetching the latest 50 URLs from the store
func BenchmarkGetCountURLs(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	// Populate the store with 1000 URLs
	newStore(10000)

	b.ResetTimer()

	// Benchmark fetching the latest 50 URLs
	for i := 0; i < b.N; i++ {
		Filter(10, "count")
	}
}

// Benchmark to test fetching the top 50 URLs based on their count
func BenchmarkGetTopURLs(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	// Populate the store with 1000 URLs and update their count
	newStore(10000)
	for i := 0; i < 1000; i++ {
		url := fmt.Sprintf("http://example%d", i)
		Update(url, true, time.Now().UnixNano())
	}

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	b.ResetTimer()

	// Benchmark fetching the top 50 URLs based on count
	for i := 0; i < b.N; i++ {
		Filter(10, "")
	}
}

// Benchmark for updating an existing URL
func BenchmarkUpdateExistingURL(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}

	// Populate the store with 1000 URLs
	newStore(10000)

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Benchmark updating an existing URL
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("http://example%d", i)
		Update(url, true, time.Now().UnixNano())
	}
}

// Benchmark add new urls
func BenchmarkAddNewURL(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}

	// Populate the store with 1000 URLs
	newStore(0)

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Benchmark updating an existing URL
	urls := make([]string, 0)
	for i := 0; i < b.N; i++ {
		urls = append(urls, fmt.Sprintf("http://example%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Update(urls[i], true, time.Now().UnixNano())
	}

}
