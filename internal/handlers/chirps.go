package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/database"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/google/uuid"
)

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func HandleCreateChirp(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{"Missing or malformed token"})
			return
		}

		userID, err := auth.ValidateJWT(tokenStr, s.JWTSecret)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{"Invalid or expired token"})
			return
		}

		var req struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"Invalid JSON"})
			return
		}

		if len(req.Body) > 140 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"Chirp is too long"})
			return
		}

		cleaned := filterProfanity(req.Body)
		now := time.Now()
		chirp, err := s.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Body:      cleaned,
			UserID:    userID,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Failed to create chirp"})
			return
		}

		writeJSON(w, http.StatusCreated, chirpResponse{
			ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID,
		})
	}
}

func HandleGetChirps(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := s.Queries.GetChirps(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Failed to fetch chirps"})
			return
		}

		resp := make([]chirpResponse, len(chirps))
		for i, chirp := range chirps {
			resp[i] = chirpResponse{
				ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID,
			}
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

func HandleGetChirp(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpIDStr := r.PathValue("chirpID")
		chirpID, err := uuid.Parse(chirpIDStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"Invalid chirp ID"})
			return
		}

		chirp, err := s.Queries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSON(w, http.StatusNotFound, ErrorResponse{"Chirp not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Failed to fetch chirp"})
			}
			return
		}

		writeJSON(w, http.StatusOK, chirpResponse{
			ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID,
		})
	}
}
