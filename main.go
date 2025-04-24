package main

import (
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	apiCfg := &apiConfig{}

	router := setupRouter(apiCfg)
	startServer(router)
}
