package main

import (
	"flag"
	"log"
	"net/http"
	"spamhaus/api"
	"spamhaus/downloader"
	"spamhaus/store"
)

func main() {
	var httpPort int
	flag.IntVar(&httpPort, "port", 80, "http port")

	flag.Parse()

	http.Handle("/submiturl", http.HandlerFunc(api.SubmitURL))

	downloader.New(3, 3)
	store.New()

	log.Printf("starting http server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("error: starting http server: %s", err.Error())
	}

}
