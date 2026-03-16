package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sudovishal/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	key := os.Getenv("POLKA_KEY")

	apiKeyHeader, err := auth.GetAPIKey(r.Header)
	if key != apiKeyHeader {
		respondWithError(w, http.StatusUnauthorized, "API Key not matching", err)
		return
	}

	type WebhookReqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := WebhookReqBody{}
	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode parameter", err)
		return
	}

	// userId, err := uuid.Parse(params.Data.UserID)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Unable to parse User ID", err)
	// }

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if params.Event == "user.upgraded" {
		_, err := cfg.db.UpdateRedMembership(r.Context(), params.Data.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "User not found", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Error while updating the membership", err)
			return
		}
		w.WriteHeader(204)

	}

}
