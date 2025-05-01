package handlers

import (
	"net/http"

	"github.com/DryHop2/chirpy/internal/state"
)

func MetricsMiddleware(s *state.State) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Metrics.Add(1)
			next.ServeHTTP(w, r)
		})
	}
}
