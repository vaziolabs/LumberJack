package internal

import (
	"context"
	"fmt"
	"net/http"
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
