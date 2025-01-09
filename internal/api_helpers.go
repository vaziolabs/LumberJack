package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"github.com/vaziolabs/lumberjack/internal/core"
	"github.com/vaziolabs/lumberjack/types"
)

// compares two byte slices for equality
func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Update the auth middleware to handle user_id from token claims
func (server *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return server.jwtConfig.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*TokenClaims)
		if !ok || claims.TokenType != "session" {
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// getNodeFromPath traverses the forest to find a node by its path
func (server *Server) getNodeFromPath(path string) (*core.Node, error) {
	if path == "" {
		return server.forest, nil
	}

	parts := strings.Split(path, "/")
	current := server.forest

	for _, part := range parts {
		found := false
		for _, child := range current.Children {
			if child.Name == part {
				current = child
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("node not found: %s", path)
		}
	}

	return current, nil
}

// UpdateSettings updates server configuration parameters
func (server *Server) UpdateSettings(userID string, settings types.ServerConfig) error {
	// Update user-specific settings
	for i := range server.forest.Users {
		if server.forest.Users[i].ID == userID {
			server.forest.Users[i].Organization = settings.Organization
			break
		}
	}

	// Update server settings if values are provided
	if settings.ServerPort != "" {
		server.config.ServerPort = settings.ServerPort
	}
	if settings.DashboardPort != "" {
		server.config.DashboardPort = settings.DashboardPort
	}
	if settings.ServerURL != "" {
		server.config.ServerURL = settings.ServerURL
	}
	if settings.DatabasePath != "" {
		server.config.DatabasePath = settings.DatabasePath
	}
	if settings.LogDirectory != "" {
		server.config.LogDirectory = settings.LogDirectory
	}
	if settings.Organization != "" {
		server.config.Organization = settings.Organization
	}
	if settings.Phone != "" {
		server.config.Phone = settings.Phone
	}

	// Save updated configuration
	return server.saveConfig()

	// TODO: Trigger a refresh of the dashboard and a reload of the server if needed
}

// saveConfig writes the current configuration to disk
func (server *Server) saveConfig() error {
	viper.Set("databases."+server.config.DatabaseName, server.config)
	return viper.WriteConfig()
}

func (server *Server) generateTokenPair(user *core.User) (*TokenPair, error) {
	// Generate session token (short-lived)
	sessionClaims := TokenClaims{
		UserID:    user.ID,
		Username:  user.Username,
		TokenType: "session",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	sessionToken := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims)
	sessionTokenString, err := sessionToken.SignedString(server.jwtConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (long-lived)
	refreshClaims := TokenClaims{
		UserID:    user.ID,
		Username:  user.Username,
		TokenType: "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(server.jwtConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		SessionToken: sessionTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (server *Server) getPaginatedLogs(page int) ([]LogEntry, bool) {
	server.logCache.mutex.RLock()
	defer server.logCache.mutex.RUnlock()

	start := (page - 1) * server.logCache.PageSize
	end := start + server.logCache.PageSize

	if start >= len(server.logCache.Logs) {
		return []LogEntry{}, false
	}

	if end > len(server.logCache.Logs) {
		end = len(server.logCache.Logs)
	}

	hasMore := end < len(server.logCache.Logs)
	return server.logCache.Logs[start:end], hasMore
}

func (server *Server) updateLogCache(level string) error {
	server.logCache.mutex.Lock()
	defer server.logCache.mutex.Unlock()

	logPath := filepath.Join(server.config.LogDirectory, fmt.Sprintf("%s.log", server.config.ProcessInfo.ID))
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		return err
	}

	// Only read new content since last update
	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// If cache exists, seek to last read position
	if server.logCache.LastOffset > 0 {
		file.Seek(server.logCache.LastOffset, 0)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := server.parseLogEntry(scanner.Text(), level)
		if err != nil {
			continue // Skip invalid entries
		}
		if entry != nil { // nil means filtered out by level
			server.logCache.Logs = append(server.logCache.Logs, *entry)
		}
	}

	// Update cache metadata
	server.logCache.LastOffset, _ = file.Seek(0, io.SeekCurrent)
	server.logCache.LastModTime = fileInfo.ModTime()

	return scanner.Err()
}

// Initialize cache only when needed
func (server *Server) initLogCacheIfNeeded() {
	if server.logCache == nil {
		server.logCache = &LogCache{
			PageSize: 100,
			Logs:     make([]LogEntry, 0),
		}
	}
}

func (server *Server) parseLogEntry(line string, level string) (*LogEntry, error) {
	// Skip empty lines
	if len(line) < 19 {
		return nil, fmt.Errorf("line too short")
	}

	// Parse timestamp
	timestamp, err := time.Parse("2006/01/02 15:04:05", line[:19])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %v", err)
	}

	// Get remainder after timestamp
	remainder := strings.TrimSpace(line[19:])

	// Parse indentation level and message
	indentCount := 0
	for _, char := range remainder {
		if char == '│' || char == '└' || char == '┌' || char == '─' {
			indentCount++
		} else {
			break
		}
	}

	// Remove tree characters and spaces
	message := strings.TrimLeft(remainder, "│└┌─ ")

	// Parse log level and message
	var logLevel string
	var content string
	var entryType string

	// Check for BEGIN/END messages
	if strings.HasPrefix(message, "BEGIN:") {
		logLevel = "INFO"
		content = strings.TrimSpace(strings.TrimPrefix(message, "BEGIN:"))
		entryType = "begin"
	} else if strings.HasPrefix(message, "END:") {
		logLevel = "INFO"
		content = strings.TrimSpace(strings.TrimPrefix(message, "END:"))
		entryType = "end"
	} else {
		// Map of symbols to log levels
		levelMap := map[string]string{
			"ℹ": "INFO",
			"✓": "SUCCESS",
			"✗": "FAILURE",
			"🔍": "DEBUG",
			"📝": "NOTICE",
			"⚠": "WARNING",
			"❌": "ERROR",
			"🔥": "CRITICAL",
			"🚨": "ALERT",
			"💀": "EMERGENCY",
		}

		// Check for level symbols
		found := false
		for symbol, lvl := range levelMap {
			if strings.HasPrefix(message, symbol) {
				logLevel = lvl
				content = strings.TrimSpace(strings.TrimPrefix(message, symbol))
				found = true
				break
			}
		}

		if !found {
			logLevel = "INFO"
			content = message
		}
		entryType = "message"
	}

	// Filter by level if specified
	if level != "" && !strings.EqualFold(level, logLevel) {
		return nil, nil
	}

	return &LogEntry{
		Timestamp: timestamp,
		Level:     logLevel,
		Message:   content,
		Type:      entryType,
		Indent:    indentCount,
	}, nil
}
