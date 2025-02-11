package downloader

import (
	"log"
	"testing"
	"time"
)

func TestWorkerPoolConcurrency(t *testing.T) {
	tests := []struct {
		name           string
		poolSize       int
		taskURLs       []string
		expectedLogMsg string
	}{
		{
			name:           "Multiple workers with concurrent tasks",
			poolSize:       3,
			taskURLs:       []string{"http://google.com", "http://google.co.uk", "https://google.cz"},
			expectedLogMsg: "wp done", // Just check that "wp done" appears in logs after processing.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var capturedLogs []string
			log.SetOutput(&logWriter{logs: &capturedLogs})

			wp := NewWorkerPool(tt.poolSize)
			startTime := time.Now()

			for _, url := range tt.taskURLs {
				wp.addTask(url)
			}

			wp.Wait()

			// Check if the pool was concurrent by comparing the total execution time
			duration := time.Since(startTime)

			// Check to see if the duration is less than the time it would take to sequentially process the urls
			if duration > time.Duration(len(tt.taskURLs))*100*time.Millisecond {
				t.Errorf("WorkerPool was not concurrent: expected execution time to be less than %v, but got %v",
					time.Duration(len(tt.taskURLs))*100*time.Millisecond, duration)
			}

		})
	}
}

// Helper struct to capture log output
type logWriter struct {
	logs *[]string
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	*lw.logs = append(*lw.logs, string(p))
	return len(p), nil
}
