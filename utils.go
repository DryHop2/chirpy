package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writePlainText(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func filterProfanity(input string) string {
	words := strings.Split(input, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if _, exists := profaneWords[lower]; exists {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
