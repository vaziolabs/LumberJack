package dashboard

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	ReadPermission Permission = iota
	WritePermission
	AdminPermission
)

func NewDashboard(apiEndpoint string, port string) *DashboardServer {
	router := mux.NewRouter()

	dashboardServer := &DashboardServer{
		apiEndpoint: apiEndpoint,
		server: &http.Server{
			Handler:      router,
			Addr:         ":" + port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("dashboard/static"))))

	// Auth routes
	router.HandleFunc("/login", dashboardServer.handleLogin).Methods("POST")

	// Protected API routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(dashboardServer.authMiddleware)
	protected.HandleFunc("/tree", dashboardServer.handleGetTree).Methods("GET")
	protected.HandleFunc("/events", dashboardServer.handleGetEvents).Methods("GET")
	protected.HandleFunc("/logs", dashboardServer.handleGetLogs).Methods("GET")
	protected.HandleFunc("/users", dashboardServer.handleGetUsers).Methods("GET")
	protected.HandleFunc("/users", dashboardServer.handleCreateUser).Methods("POST")

	// Main dashboard routes
	router.HandleFunc("/", dashboardServer.handleLoginPage).Methods("GET")
	router.Handle("/dashboard", dashboardServer.authMiddleware(http.HandlerFunc(dashboardServer.handleDashboard))).Methods("GET")

	return dashboardServer
}

func (s *DashboardServer) Start() error {
	go func() {
		log.Printf("%s", "Dashboard server starting on http://localhost:"+s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Dashboard server error: %v", err)
		}
	}()
	return nil
}

func (s *DashboardServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
