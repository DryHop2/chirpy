package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/database"
	"github.com/DryHop2/chirpy/internal/state"
)

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleUpdateUser(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.ValidateJWTFromRequest(r, s.JWTSecret)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid or missing token"})
			return
		}

		var req updateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
			return
		}

		err = s.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to update user"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"email": req.Email,
		})
	}
}
