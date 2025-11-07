package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Println("Server listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
		log.Fatal(err)
	}
}