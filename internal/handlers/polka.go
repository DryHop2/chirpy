package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func HandlePolkaWebhook(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != s.PolkaKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req polkaWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		if req.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = s.Queries.UpgradeUserToChirpyRed(r.Context(), req.Data.UserID)
		if err != nil {
			// This would loop forever
			// writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "User not found"})
			log.Printf("Polka webhook: user not found: %s", req.Data.UserID)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
