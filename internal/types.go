package internal

import (
	"net/http"
	"sync"
	"time"

	"github.com/vaziolabs/LumberJack/internal/core"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey []byte
	ExpiresIn time.Duration
}

type Server struct {
	forest    *core.Node
	mutex     sync.Mutex
	jwtConfig JWTConfig
	logger    Logger
	server    *http.Server
}
