package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))))

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Println("listening on http://localhost:8080")
	log.Println("  - health:  GET /healthz")
	log.Println("  - app:     GET /app")
	log.Println("  - assets:  GET /app/assets/")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}