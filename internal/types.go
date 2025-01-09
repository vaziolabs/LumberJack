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

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Trace     string    `json:"trace,omitempty"`
	Indent    int       `json:"indent"`
	Type      string    `json:"type"`
}

type LogCache struct {
	Logs        []LogEntry
	LastOffset  int64
	LastModTime time.Time
	PageSize    int
	mutex       sync.RWMutex
}

type Server struct {
	forest    *core.Node
	mutex     sync.Mutex
	jwtConfig JWTConfig
	logger    types.Logger
	server    *http.Server
	config    types.ServerConfig
	lastHash  []byte
	logCache  *LogCache
}
