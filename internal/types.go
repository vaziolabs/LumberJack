package internal

import (
	"net/http"
	"sync"
	"time"

	"github.com/vaziolabs/lumberjack/internal/core"
	"github.com/vaziolabs/lumberjack/types"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SessionKey []byte
	RefreshKey []byte
	ExpiresIn  time.Duration
	SecretKey  []byte
}

type Server struct {
	forest    *core.Node
	mutex     sync.Mutex
	jwtConfig JWTConfig
	logger    types.Logger
	server    *http.Server
	config    types.ServerConfig
	lastHash  []byte
}
