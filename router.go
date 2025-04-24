package main

import (
	"net/http"
)

func setupRouter(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	// mux.Handle("/", http.FileServer(http.Dir(".")))
	appFs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.middlewareMetricsInc(appFs))

	// assetsHandler := http.FileServer(http.Dir("./assets"))
	// mux.Handle("assets/", assetsHandler)

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /api/reset", cfg.handleReset)

	return mux
}
