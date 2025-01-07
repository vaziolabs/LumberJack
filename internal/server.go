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

	server.logger.Enter("NewServer")
	defer server.logger.Exit("NewServer")

	// Load existing database if it exists
	dbPath := filepath.Join(config.DbPath, config.DbName+".dat")
	if err := server.loadFromFile(dbPath); err == nil {
		server.logger.Info("Loaded existing database from %s", dbPath)
	} else if config.User.Username != "" && config.User.Password != "" {
		server.logger.Info("Creating new database")

		adminUser := core.User{
			ID:       core.GenerateID(),
			Username: config.User.Username,
			Email:    config.User.Email,
		}

		server.logger.Info("Creating admin user with username: %s", config.User.Username)

		if err := adminUser.SetPassword(config.User.Password); err != nil {
			server.logger.Failure("failed to set admin password: %v", err)
			return nil, err
		}
		server.logger.Debug("Admin user: %+v", adminUser)

		server.logger.Info("Admin password set successfully")

		if err := server.forest.AssignUser(adminUser, core.AdminPermission); err != nil {
			server.logger.Failure("failed to save admin user: %v", err)
			return nil, err
		}

		if err := server.writeChangesToFile(server.forest, dbPath); err != nil {
			server.logger.Failure("failed to save state after user creation: %v", err)
			return nil, err
		}

		server.logger.Info("Database created successfully with admin user: %s", adminUser.Username)
	} else {
		server.logger.Failure("failed to create database: %v", err)
		return nil, err
	}

	return server, nil
}

func (s *Server) Start() error {
	if s.server == nil {
		return errors.New("server not initialized")
	}

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", s.handleLogin).Methods("POST")
	router.HandleFunc("/create_user", s.handleCreateUser).Methods("POST")
	// Protected routes
	router.HandleFunc("/end_event", s.authMiddleware(s.handleEndEvent)).Methods("POST")
	router.HandleFunc("/append_event", s.authMiddleware(s.handleAppendToEvent)).Methods("POST")
	router.HandleFunc("/start_event", s.authMiddleware(s.handleStartEvent)).Methods("POST")
	router.HandleFunc("/get_event_entries", s.authMiddleware(s.handleGetEventEntries)).Methods("POST")
	router.HandleFunc("/plan_event", s.authMiddleware(s.handlePlanEvent)).Methods("POST")
	router.HandleFunc("/assign_user", s.authMiddleware(s.handleAssignUser)).Methods("POST")
	router.HandleFunc("/start_time_tracking", s.authMiddleware(s.handleStartTimeTracking)).Methods("POST")
	router.HandleFunc("/stop_time_tracking", s.authMiddleware(s.handleStopTimeTracking)).Methods("POST")
	router.HandleFunc("/get_time_tracking", s.authMiddleware(s.handleGetTimeTracking)).Methods("GET")
	router.HandleFunc("/get_tree", s.authMiddleware(s.handleGetTree)).Methods("GET")
	router.HandleFunc("/get_users", s.authMiddleware(s.handleGetUsers)).Methods("GET")

	s.server.Handler = router
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
