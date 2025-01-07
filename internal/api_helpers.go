package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

// Add middleware for authentication
func (server *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return server.jwtConfig.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
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
}

// saveConfig writes the current configuration to disk
func (server *Server) saveConfig() error {
	viper.Set("databases."+server.config.DatabaseName, server.config)
	return viper.WriteConfig()
}
