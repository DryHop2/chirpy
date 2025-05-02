package main

import (
	"net/http"

	"github.com/DryHop2/chirpy/internal/handlers"
	"github.com/DryHop2/chirpy/internal/state"
)

func setupRouter(s *state.State) *http.ServeMux {
	mux := http.NewServeMux()

	appFs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", handlers.MetricsMiddleware(s)(appFs))

	mux.HandleFunc("GET /api/healthz", handlers.HandleReadiness)
	mux.HandleFunc("POST /api/users", handlers.HandleCreateUser(s))
	mux.HandleFunc("PUT /api/users", handlers.HandleUpdateUser(s))
	mux.HandleFunc("POST /api/login", handlers.HandleLogin(s))
	mux.HandleFunc("POST /api/refresh", handlers.HandleRefresh(s))
	mux.HandleFunc("POST /api/revoke", handlers.HandleRevoke(s))
	mux.HandleFunc("POST /api/chirps", handlers.HandleCreateChirp(s))
	mux.HandleFunc("POST /api/polka/webhooks", handlers.HandlePolkaWebhook(s))
	mux.HandleFunc("GET /api/chirps", handlers.HandleGetChirps(s))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.HandleGetChirp(s))
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", handlers.HandleDeleteChirps(s))

	mux.HandleFunc("POST /admin/reset", handlers.HandleAdminReset(s))
	mux.HandleFunc("GET /admin/metrics", handlers.HandleAdminMetrics(s))

	return mux
}
