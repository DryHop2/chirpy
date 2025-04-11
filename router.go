package main

import (
	"net/http"
)

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))

	mux.Handle("/", fs)
	return mux
}
