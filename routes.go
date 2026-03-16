package main

import (
	"net/http"
)

func (apiCfg *apiConfig) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset) // delete All Users

	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerCreateRFToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpbyID)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateCreds)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)
}
