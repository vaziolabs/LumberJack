package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vaziolabs/lumberjack/cmd"
	"github.com/vaziolabs/lumberjack/internal/core"
)

func NewServer(port string, config cmd.User) *Server {
	router := mux.NewRouter()

	server := &Server{
		forest: core.NewForest("root"),
		jwtConfig: JWTConfig{
			SecretKey: []byte("your-secret-key"),
			ExpiresIn: 24 * time.Hour,
		},
		logger: NewLogger(),
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}

	// Create admin user from config
	adminUser := core.User{
		Username: config.Username,
		Password: hashPassword(config.Password),
	}

	if err := server.forest.AssignUser(adminUser, core.AdminPermission); err != nil {
		server.logger.Failure("Failed to save admin user: %v", err)
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

	return server
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("API server starting on http://localhost:" + s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Failure("API server error: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
