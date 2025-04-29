package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/DryHop2/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	DB             *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		DB:             dbQueries,
	}

	router := setupRouter(apiCfg)
	startServer(router)
}
