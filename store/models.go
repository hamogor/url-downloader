package store

import "time"

type URLData struct {
	Count          int       // number of times the url has been submitted
	Successes      int       // number of successful downloads
	Failures       int       // number of failed downloads
	LastDownloadMs int64     // Last download duration (ms)
	LastSubmitted  time.Time // Time of last download
}
