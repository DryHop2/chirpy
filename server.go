package main

import (
	"log"
	"net/http"
)

func startServer(handler http.Handler) {
	addr := ":8080"
	log.Printf("Starting server on %s...", addr)
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
