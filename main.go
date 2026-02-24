package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/sudovishal/chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	jwtSecret      string
}

type cleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type parameter struct {
	Body string `json:"body"`
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error opening database %s", err)
	}

	dbQueries := database.New(db)
	secretStr := os.Getenv("jwtsecret")
	apiCfg := apiConfig{
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
		jwtSecret:      secretStr,
	}
	mux := http.NewServeMux()

	apiCfg.registerRoutes(mux)

	server := http.Server{Addr: ":8080", Handler: mux}

	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 	// w.WriteHeader(http.StatusMethodNotAllowed)
	// 	w.Write([]byte("OK"))
	// })

	fmt.Println("Server listening on localhost:8080")
	log.Fatal(server.ListenAndServe())

}
