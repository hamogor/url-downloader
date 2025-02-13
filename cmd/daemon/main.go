package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spamhaus/api"
	"spamhaus/downloader"
	"spamhaus/store"
	"syscall"
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

	store.New()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go server.ListenAndServe()

	httpServer, err := api.Start(config.Server.Port)
	if err != nil {
		log.Fatalf("error starting http server: %v", err)
	}

	downloader.NewBatchProcess(
		time.Duration(config.Downloader.BatchIntervalSeconds),
		config.Downloader.WorkerPoolSize,
		config.Downloader.NumOfBatchURLs,
	)

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)

	<-shutdown
	api.Shutdown(httpServer)
	store.Shutdown()

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
