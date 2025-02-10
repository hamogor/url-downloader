package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"spamhaus/downloader"
)

type SubmitURLRequest struct {
	URL string `json:"url"`
}

func SubmitURL(w http.ResponseWriter, r *http.Request) {

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
	go downloader.GlobalDownloader.DownloadURl(req.URL)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "url submitted"})

}

func TopURLs(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	sortBy := r.URL.Query().Get("sort_by")

	if sortBy != "count" && sortBy != "requests" {
		http.Error(w, fmt.Sprintf("error: invalid sort by %s", sortBy), http.StatusBadRequest)
	}

	if orderBy != "asc" && orderBy != "desc" {
		http.Error(w, fmt.Sprintf("error: invalid sort by %s", orderBy), http.StatusBadRequest)
	}

	// return a json response
	/*
		[
			{"url": "https://example.com", count: 1},
		]
	*/
}

// validateURL checks if a given URL is valid.
func (u *SubmitURLRequest) isValidURL() error {
	parsedURL, err := url.ParseRequestURI(u.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid URL format")
	}
	return nil
}
