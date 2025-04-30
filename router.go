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
	mux.HandleFunc("GET /admin/metrics", cfg.handleAdminMetrics)
	// mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	// mux.HandleFunc("POST /api/validate_chirp", cfg.handleValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handleCreateUser)
	mux.HandleFunc("POST /admin/reset", cfg.handleAdminReset)
	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChrip)

	return mux
}
