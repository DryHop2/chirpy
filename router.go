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

	mux.HandleFunc("GET /healthz", handleReadiness)
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /reset", cfg.handleReset)

	return mux
}
