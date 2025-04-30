package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DryHop2/chirpy/internal/database"
	"github.com/google/uuid"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	writePlainText(w, http.StatusOK, "OK")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHits.Load()

	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, hits)

	writePlainText(w, http.StatusOK, html)
}

// func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
// 	cfg.fileServerHits.Store(0)
// 	writePlainText(w, http.StatusOK, "Counter reset.")
// }

type ChirpRequest struct {
	Body string `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var chirp ChirpRequest
	err := decoder.Decode(&chirp)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "Something went wrong",
		})
		return
	}

	if len(chirp.Body) > 140 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "Chirp is too long",
		})
		return
	}

	cleaned := filterProfanity(chirp.Body)

	writeJSON(w, http.StatusOK, CleanedResponse{
		CleanedBody: cleaned,
	})
}

type createUserRequest struct {
	Email string `json:"email"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), req.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Could not create user"})
		return
	}

	resp := userResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (cfg *apiConfig) handleAdminReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to reset users", http.StatusInternalServerError)
		return
	}

	writePlainText(w, http.StatusOK, "All users deleted.")
}

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChrip(w http.ResponseWriter, r *http.Request) {
	var req createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	if len(req.Body) > 140 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Chirp is too long"})
		return
	}

	cleaned := filterProfanity(req.Body)
	now := time.Now()

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      cleaned,
		UserID:    req.UserID,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to create chirp"})
		return
	}

	writeJSON(w, http.StatusCreated, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
