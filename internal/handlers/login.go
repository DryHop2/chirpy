package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/database"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

		token, err := auth.MakeJWT(dbUser.ID, s.JWTSecret, time.Hour)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Token generation failed"})
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to create refresh token"})
			return
		}

		now := time.Now().UTC()
		expiresAt := now.Add(60 * 24 * time.Hour)

		err = s.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    dbUser.ID,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to store refresh token"})
			return
		}

		writeJSON(w, http.StatusOK, struct {
			ID           uuid.UUID `json:"id"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Email        string    `json:"email"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
			IsChirpyRed  bool      `json:"is_chirpy_red"`
		}{
			ID:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			Token:        token,
			RefreshToken: refreshToken,
			IsChirpyRed:  dbUser.IsChirpyRed,
		})
	}
}
