package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/DryHop2/chirpy/internal/database"
	"github.com/DryHop2/chirpy/internal/state"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "google.golang.org/grpc/balancer/grpclb/state"
)

func main() {
	err := godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	dbQueries := database.New(db)

	appState := &state.State{
		Queries:   dbQueries,
		Platform:  os.Getenv("PLATFORM"),
		JWTSecret: jwtSecret,
	}

	router := setupRouter(appState)
	startServer(router)
}
