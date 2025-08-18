package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	currentDir := "."
	port := "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	handlerFile := http.StripPrefix("/app", http.FileServer(http.Dir(currentDir)))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFile))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidator)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
