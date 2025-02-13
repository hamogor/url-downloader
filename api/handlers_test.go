package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spamhaus/store"
	"testing"
)

// Mock the store.Filter function to return dummy data for testing
func newStore() {
	store.New()
	for i := 0; i < 2; i++ {
		store.Update(
			fmt.Sprintf("http://example%d.com", i),
			true,
			int64(100+i),
		)
	}
}
func TestSubmitURL(t *testing.T) {

	tests := []struct {
		name            string
		body            string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:           "valid URL submission",
			body:           `{"url": "http://validurl.com"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid URL format",
			body:           `{"url": "invalid-url"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with the provided body
			req, err := http.NewRequest(http.MethodPost, "/submit-url", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the SubmitURL handler
			handler := http.HandlerFunc(SubmitURL)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}

		})
	}
}

func TestTopURLs(t *testing.T) {
	// Prepare the test data
	store.Update("http://example0.com", true, 100)
	store.Update("http://example1.com", true, 200)
	store.Update("http://example2.com", true, 150)

	tests := []struct {
		name             string
		sortBy           string
		getTopN          string
		expectedStatus   int
		expectedResponse []TopURLSResponse
	}{
		{
			name:           "valid request for top URLs sorted by count",
			sortBy:         "count",
			getTopN:        "2",
			expectedStatus: http.StatusOK,
			expectedResponse: []TopURLSResponse{
				{URL: "http://example1.com", Count: 200},
				{URL: "http://example2.com", Count: 150},
			},
		},
		{
			name:           "valid request for top URLs sorted by latest",
			sortBy:         "latest",
			getTopN:        "2",
			expectedStatus: http.StatusOK,
			expectedResponse: []TopURLSResponse{
				{URL: "http://example2.com", Count: 150}, // This would depend on the latest data logic
				{URL: "http://example1.com", Count: 200},
			},
		},
		{
			name:           "invalid sort by parameter",
			sortBy:         "invalid",
			getTopN:        "2",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid n parameter",
			sortBy:         "count",
			getTopN:        "not-a-number",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty getTopN value",
			sortBy:         "count",
			getTopN:        "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with query parameters
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/top-urls?sort_by=%s&get_n=%s", tt.sortBy, tt.getTopN), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the TopURLs handler
			handler := http.HandlerFunc(TopURLs)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}

			// Check the response body for valid results
			if rr.Code == http.StatusOK {
				var res []TopURLSResponse
				err := json.Unmarshal(rr.Body.Bytes(), &res)
				if err != nil {
					t.Fatalf("could not unmarshal response: %v", err)
				}

				if len(res) != len(tt.expectedResponse) {
					t.Fatalf("expected response length %d, got %d", len(tt.expectedResponse), len(res))
				}

				// Check each returned URL in response
				for i, item := range res {
					if item.URL != tt.expectedResponse[i].URL || item.Count != tt.expectedResponse[i].Count {
						t.Errorf("expected %v, got %v", tt.expectedResponse[i], item)
					}
				}
			}
		})
	}
}
