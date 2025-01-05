package internal

import (
	"forestree"
	"logger"
	"sync"
	"time"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey []byte
	ExpiresIn time.Duration
}

// App represents the main application structure
type App struct {
	forest    *forestree.Node
	mutex     sync.Mutex
	jwtConfig JWTConfig
	logger    logger.Logger
}
