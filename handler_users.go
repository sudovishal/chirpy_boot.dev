package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sudovishal/chirpy/internal/auth"
	"github.com/sudovishal/chirpy/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	// fmt.Println(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode parameter", err)
		// w.WriteHeader(400)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		// w.WriteHeader(500)
		return
	}

	user, err := cfg.db.CreateUser(
		r.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPwd,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// apiUser := User{
	// 	ID:             user.ID,
	// 	CreatedAt:      user.UpdatedAt.Time,
	// 	UpdatedAt:      user.UpdatedAt.Time,
	// 	Email:          user.Email.String,
	// 	HashedPassword: user.HashedPassword,
	// }
	// // fmt.Println(apiUser.ID)

	// w.WriteHeader(201)
	// json.NewEncoder(w).Encode(apiUser)
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt.Time,
			UpdatedAt:   user.UpdatedAt.Time,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(403)
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		// log.Printf("Error deleting users: %s", err)
		// w.WriteHeader(500)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "reset successful",
	}

	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, resp)
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// r.Header.Add("Content-Type", "application/json")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	// fmt.Println(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// const defaultExpiration = 3600
	// const maxExpiration = 3600

	// expiresIn := defaultExpiration

	// if params.ExpiresInSeconds != nil {
	// 	requested := *params.ExpiresInSeconds
	// 	expiresIn = min(requested, maxExpiration)
	// }

	expiresAt := time.Now().Add(1 * time.Hour)

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// apiUser := User{
	// 	ID:             user.ID,
	// 	CreatedAt:      user.CreatedAt.Time,
	// 	UpdatedAt:      user.UpdatedAt.Time,
	// 	Email:          user.Email.String,
	// 	HashedPassword: user.HashedPassword,
	// }

	// fmt.Println(params.Password, user.HashedPassword)
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Incorrect Email or password", err)
	// 	// log.Printf("Error comparing passwords: %s", err)
	// 	// w.WriteHeader(401)
	// 	return
	// }

	// if !passwordVerify {
	// 	w.WriteHeader(401)
	// 	return
	// } else {
	// 	w.Header().Add("Content-Type", "application/json")
	// 	w.WriteHeader(200)
	// 	json.NewEncoder(w).Encode(apiUser)
	// }

	// fmt.Println(time.Until(expiresAt))
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Until(expiresAt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	rfToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate Refresh Token String", err)
		return
	}

	// expirationRF := time.Now().Add(60 * 24 * time.Hour)
	refreshToken, err := cfg.db.GenerateRefreshToken(r.Context(), database.GenerateRefreshTokenParams{
		Token:     rfToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store Refresh Token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt.Time,
			UpdatedAt:   user.UpdatedAt.Time,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})

}

func (cfg *apiConfig) handlerCreateRFToken(w http.ResponseWriter, r *http.Request) {

	tokenHeader, err := auth.GetBearerToken(r.Header)
	fmt.Println(tokenHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch Bearer Token", err)
		return
	}

	User, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to fetch user", err)
		return
	}

	// expiresAt := time.Now().Add(1 * time.Hour)
	refreshAccessToken, err := auth.MakeJWT(User.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: refreshAccessToken,
	})

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenHeader, err := auth.GetBearerToken(r.Header)
	fmt.Println(tokenHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch Bearer Token", err)
		return
	}

	err = cfg.db.DeleteRFToken(r.Context(), tokenHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpdateCreds(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unable to fetch Bearer Token", err)
		return
	}

	type request struct {
		Email           string `json:"email"`
		UpdatedPassword string `json:"updated_password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := request{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error validating token:", err)
		return
	}

	hashedPwd, err := auth.HashPassword(params.UpdatedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.UpdateEmailPass(r.Context(), database.UpdateEmailPassParams{
		Email:          params.Email,
		ID:             userId,
		HashedPassword: hashedPwd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update the Email and Password", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userId,
			Email:     params.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		},
	})
}
