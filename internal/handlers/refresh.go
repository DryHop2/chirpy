package handlers

import (
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/state"
)

func HandleRefresh(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Missing or malformed token"})
			return
		}

		user, err := s.Queries.GetUserFromRefreshToken(r.Context(), tokenStr)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid or expired refresh token"})
			return
		}

		newToken, err := auth.MakeJWT(user.ID, s.JWTSecret, time.Hour)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate access token"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"token": newToken,
		})
	}
}
