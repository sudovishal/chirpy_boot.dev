package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sudovishal/chirpy/internal/auth"
	"github.com/sudovishal/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type parameter struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	fmt.Println(params.Body)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	censored := removeProfane(params.Body)

	cleanedResponse := cleanResponse{
		CleanedBody: censored,
	}

	if len(params.Body) > 140 {
		response := errorResponse{Error: "Chirp is too long"}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	} else if params.Body == "" {
		response := errorResponse{Error: "Body is required"}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	// response := validResponse{Valid: true}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(cleanedResponse)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	type reqPayload struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqPayload{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	if len(params.Body) > 140 {
		response := errorResponse{Error: "Chirp is too long"}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	} else if params.Body == "" {
		response := errorResponse{Error: "Body is required"}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		w.WriteHeader(400)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error validating token:", err)
		// log.Printf("Error validating token: %s", err)
		// w.WriteHeader(401)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating chirp:", err)
		// log.Printf("Error creating chirp: %s", err)
		// w.WriteHeader(500)
		return
	}

	resChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resChirp)

}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		// w.WriteHeader(500)
		return
	}

	resChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		resChirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resChirps); err != nil {
		log.Printf("Error encoding chirps response: %s", err)
	}
	// json.NewEncoder(w).Encode(resChirps)
}

func (cfg *apiConfig) getChirpbyID(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")

	chirpID := r.PathValue("chirpID")

	if chirpID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "chirpID is required"})
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Fatalf("failed to parse UUID: %v", err)
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "chirp not found"})
		return
	}

	resChirp := Chirp{
		ID:        chirpUUID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resChirp)
}
