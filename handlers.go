package main

import (
	"fmt"
	"net/http"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	writePlainText(w, http.StatusOK, "Ok")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHits.Load()
	writePlainText(w, http.StatusOK, fmt.Sprintf("Hits: %d", hits))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	writePlainText(w, http.StatusOK, "Counter reset.")
}
