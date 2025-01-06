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
			Username: config.User.Username,
			Password: hashPassword(config.User.Password),
		}

		if err := server.forest.AssignUser(adminUser, core.AdminPermission); err != nil {
			server.logger.Failure("failed to save admin user: %v", err)
			return nil, err
		}

		if err := server.writeChangesToFile(server.forest, dbPath); err != nil {
			server.logger.Failure("failed to save initial database state: %v", err)
			return nil, err
		} else {
			server.logger.Info("Database created successfully")
		}

		// We exit early if we are simply creating a new database
		return nil, nil
	} else {
		server.logger.Failure("failed to create database: %v", err)
		return nil, err
	}

	// Public routes
	router.HandleFunc("/login", server.handleLogin).Methods("POST")
	router.HandleFunc("/create_user", server.handleCreateUser).Methods("POST")

	// Protected routes
	router.HandleFunc("/end_event", server.authMiddleware(server.handleEndEvent)).Methods("POST")
	router.HandleFunc("/append_event", server.authMiddleware(server.handleAppendToEvent)).Methods("POST")
	router.HandleFunc("/start_event", server.authMiddleware(server.handleStartEvent)).Methods("POST")
	router.HandleFunc("/get_event_entries", server.authMiddleware(server.handleGetEventEntries)).Methods("POST")
	router.HandleFunc("/plan_event", server.authMiddleware(server.handlePlanEvent)).Methods("POST")
	router.HandleFunc("/assign_user", server.authMiddleware(server.handleAssignUser)).Methods("POST")
	router.HandleFunc("/start_time_tracking", server.authMiddleware(server.handleStartTimeTracking)).Methods("POST")
	router.HandleFunc("/stop_time_tracking", server.authMiddleware(server.handleStopTimeTracking)).Methods("POST")
	router.HandleFunc("/get_time_tracking", server.authMiddleware(server.handleGetTimeTracking)).Methods("GET")
	router.HandleFunc("/get_tree", server.authMiddleware(server.handleGetTree)).Methods("GET")
	router.HandleFunc("/get_users", server.authMiddleware(server.handleGetUsers)).Methods("GET")

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
