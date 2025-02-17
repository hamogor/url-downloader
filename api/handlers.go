package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"spamhaus/downloader"
	"spamhaus/store"
	"strconv"
)

type SubmitURLRequest struct {
	URL string `json:"url"`
}

type TopURLSResponse struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

func SubmitURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Grab and decode the request body
	var req SubmitURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: decoding url from SubmitURLRequest: %s", err), http.StatusInternalServerError)
		return
	}

	// Validate that the url given is a valid URL
	err = req.isValidURL()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: validating url from SubmitURLRequest: %s", err), http.StatusBadRequest)
		return
	}

	// Add download job for this URL to the worker pool
	go downloader.AddTask(req.URL)

	err = json.NewEncoder(w).Encode(map[string]string{"message": "url submitted"})
	if err != nil {
		http.Error(w, fmt.Sprintf("error: encoding url to SubmitURLRequest: %s", err), http.StatusInternalServerError)
		return
	}

}

func TopURLs(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Fetch and validate query params
	sortBy := r.URL.Query().Get("sort_by")
	getTopN := r.URL.Query().Get("get_n")
	if sortBy != "count" && sortBy != "latest" {
		http.Error(w, fmt.Sprintf("error: invalid sort by %s", sortBy), http.StatusBadRequest)
	}

	n, err := strconv.Atoi(getTopN)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: invalid n: %s should be convertable to int", getTopN), http.StatusBadRequest)
	}

	// Filter for the latest n URLs
	urls := store.Filter(n, sortBy)
	responses := make([]TopURLSResponse, 0, n)
	for _, node := range urls {
		responses = append(responses, TopURLSResponse{
			URL:   node.URL,
			Count: node.Data.Count,
		})
	}

	jsonData, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: encoding top urls: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: writing top urls: %s", err), http.StatusInternalServerError)
	}

}

// validateURL checks if a given URL is valid.
func (u *SubmitURLRequest) isValidURL() error {
	parsedURL, err := url.ParseRequestURI(u.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid URL format")
	}
	return nil
}
