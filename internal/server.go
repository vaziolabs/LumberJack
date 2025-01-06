package internal

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/vaziolabs/lumberjack/internal/core"
	"github.com/vaziolabs/lumberjack/types"
)

func NewServer(config types.ServerConfig) (*Server, error) {
	router := mux.NewRouter()

	server := &Server{
		forest: core.NewForest("root"),
		jwtConfig: JWTConfig{
			SecretKey: []byte("your-secret-key"),
			ExpiresIn: 24 * time.Hour,
		},
		logger: types.NewLogger(),
		server: &http.Server{
			Addr:    ":" + config.Port,
			Handler: router,
		},
	}

	// Load existing database if it exists
	dbPath := filepath.Join(config.DbPath, config.DBName+".dat")
	if err := server.readChangesFromFile(dbPath, server.forest); err == nil {
		server.logger.Info("Loaded existing database from %s", dbPath)
	} else if config.User.Username != "" && config.User.Password != "" {
		server.logger.Info("Creating new database")

		adminUser := core.User{
			ID:       core.GenerateID(),
			Username: config.User.Username,
		}

		if err := adminUser.SetPassword(config.User.Password); err != nil {
			server.logger.Failure("failed to set admin password: %v", err)
			return nil, err
		}

		if err := server.forest.AssignUser(adminUser, core.AdminPermission); err != nil {
			server.logger.Failure("failed to save admin user: %v", err)
			return nil, err
		}

		if err := server.writeChangesToFile(server.forest, dbPath); err != nil {
			server.logger.Failure("failed to save initial database state: %v", err)
			return nil, err
		}

		server.logger.Info("Database created successfully with admin user: %s", adminUser.Username)
		return server, nil
	}

	return server, nil
}

func (s *Server) Start() error {
	if s.server == nil {
		return errors.New("server not initialized")
	}

	go func() {
		s.logger.Info("API server starting on http://localhost" + s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Failure("API server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		s.logger.Info("Shutting down API server")
		return s.server.Shutdown(ctx)
	}
	return nil
}
