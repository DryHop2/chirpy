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

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func HandleCreateUser(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"Invalid JSON"})
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"Hashing failed"})
			return
		}

		dbUser, err := s.Queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"User creation failed"})
			return
		}

		writeJSON(w, http.StatusCreated, userResponse{
			ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Email: dbUser.Email,
		})
	}
}
