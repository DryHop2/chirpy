package handlers

import (
	"fmt"
	"net/http"

	"github.com/DryHop2/chirpy/internal/state"
)

func HandleAdminReset(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Platform != "dev" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if err := s.Queries.DeleteAllUsers(r.Context()); err != nil {
			http.Error(w, "Failed to reset users", http.StatusInternalServerError)
			return
		}
		writePlainText(w, http.StatusOK, "All users deleted.")
	}
}

func HandleAdminMetrics(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
			<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>`,
			s.Metrics.Load())
		writePlainText(w, http.StatusOK, html)
	}
}
