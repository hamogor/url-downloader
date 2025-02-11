package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"spamhaus/api"
	"spamhaus/downloader"
	"spamhaus/store"
	"time"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Downloader struct {
		WorkerPoolSize       int `yaml:"worker_pool_size"`
		NumOfBatchURLs       int `yaml:"num_of_batch_urls"`
		BatchIntervalSeconds int `yaml:"batch_interval_seconds"`
	} `yaml:"downloader"`
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	http.Handle("/submiturl", http.HandlerFunc(api.SubmitURL))
	http.Handle("/topurls", http.HandlerFunc(api.TopURLs))

	downloader.New(config.Downloader.WorkerPoolSize)
	store.New()

	downloader.NewBatchProcess(time.Duration(config.Downloader.BatchIntervalSeconds), config.Downloader.WorkerPoolSize, config.Downloader.NumOfBatchURLs)

	log.Printf("starting http server")
	if err := http.ListenAndServe(config.Server.Port, nil); err != nil {
		log.Printf("error: starting http server: %s", err.Error())
	}

}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
