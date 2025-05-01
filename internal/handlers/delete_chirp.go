package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/google/uuid"
)

func HandleDeleteChirps(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpIDStr := r.PathValue("chirpID")
		chirpID, err := uuid.Parse(chirpIDStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid chirp ID"})
			return
		}

		userID, err := auth.ValidateJWTFromRequest(r, s.JWTSecret)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid or missing token"})
			return
		}

		ownerID, err := s.Queries.GetChirpOwner(r.Context(), chirpID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "Chirp not found."})
			} else {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to look up chirp"})
			}
			return
		}

		if ownerID != userID {
			writeJSON(w, http.StatusForbidden, ErrorResponse{Error: "You are not allowed to delete this chirp"})
			return
		}

		err = s.Queries.DeleteChirp(r.Context(), chirpID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete chirp"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
