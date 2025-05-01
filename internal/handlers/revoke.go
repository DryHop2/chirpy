package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/database"
	"github.com/DryHop2/chirpy/internal/state"
)

func HandleRevoke(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Missing or malformed token"})
			return
		}

		now := time.Now().UTC()

		err = s.Queries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
			Token: tokenStr,
			RevokedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		})

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to revoke token."})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
