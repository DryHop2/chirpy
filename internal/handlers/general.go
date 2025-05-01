package handlers

import (
	"fmt"
	"net/http"
)

func HandleReadiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Auth header in readiness:", r.Header.Get("Authorization"))
	writePlainText(w, http.StatusOK, "OK")
}
