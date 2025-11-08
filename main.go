package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	cfg := &apiConfig{}


	normalizeMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.ReplaceAll(r.URL.Path, "//", "/")
			next.ServeHTTP(w, r)
		})
	}

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", cfg.handleAdminMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)

	fsHandler := http.FileServer(http.Dir("./assets"))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(
		http.StripPrefix("/app/assets/", fsHandler),
	))

	appHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})
	mux.Handle("/app", cfg.middlewareMetricsInc(appHandler))

	server := &http.Server{
		Addr: ":8080",
		Handler: normalizeMiddleware(mux),
	}

	log.Println("listening on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileServerHits.Load()
	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, count)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
