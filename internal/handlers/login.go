package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/google/uuid"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func HandleLogin(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"Invalid JSON"})
			return
		}

		dbUser, err := s.Queries.GetUserByEmail(r.Context(), req.Email)
		if err != nil || auth.CheckPasswordHash(dbUser.HashedPassword, req.Password) != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{"Incorrect email or password"})
			return
		}

		expiration := time.Hour
		if req.ExpiresInSeconds > 0 && req.ExpiresInSeconds <= 3600 {
			expiration = time.Duration(req.ExpiresInSeconds) * time.Second
		}

		token, err := auth.MakeJWT(dbUser.ID, s.JWTSecret, expiration)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Token generation failed"})
			return
		}

		writeJSON(w, http.StatusOK, struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}{
			ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Email: dbUser.Email, Token: token,
		})
	}
}
