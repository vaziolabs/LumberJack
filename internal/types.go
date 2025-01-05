package internal

import (
	"sync"
	"time"

	"github.com/vaziolabs/LumberJack/internal/core"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey []byte
	ExpiresIn time.Duration
}

// App represents the main application structure
type App struct {
	Forest    *core.Node
	Mutex     sync.Mutex
	JWTConfig JWTConfig
	Logger    Logger
}
