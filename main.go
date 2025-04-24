package main

import (
	"sync/atomic"
)

type apiConfig struct {
	fileSeverHits atomic.Int32
}

func main() {
	apiCfg := &apiConfig{}

	router := setupRouter(apiCfg)
	startServer(router)
}
