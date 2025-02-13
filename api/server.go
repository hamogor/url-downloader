package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(port string) (*http.Server, error) {
	router := http.NewServeMux()
	router.Handle("/submiturl", http.HandlerFunc(SubmitURL))
	router.Handle("/topurls", http.HandlerFunc(TopURLs))

	httpServer := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s", port),
	}

	go func() error {
		err := httpServer.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	}()

	log.Printf("http: now serving http server on: %s", port)

	return httpServer, nil
}

func Shutdown(server *http.Server) {
	log.Println("http: attempting graceful shutdown")
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Printf("http: failed to shutdown gracefully: %s", err)
	}
	log.Println("http: shutdown complete")
}
