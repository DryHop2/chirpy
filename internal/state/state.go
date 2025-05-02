package state

import (
	"sync/atomic"

	"github.com/DryHop2/chirpy/internal/database"
)

type State struct {
	Queries   *database.Queries
	JWTSecret string
	PolkaKey  string
	Platform  string
	Metrics   atomic.Int32
}
