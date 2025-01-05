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

func NewDashboardServer(apiEndpoint string, dashboardPort string) *DashboardServer {
	router := mux.NewRouter()

	dashboardServer := &DashboardServer{
		apiEndpoint: apiEndpoint,
		server: &http.Server{
			Handler:      router,
			Addr:         ":" + dashboardPort,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("dashboard/static"))))

	// API routes
	router.HandleFunc("/api/tree", dashboardServer.handleGetTree).Methods("GET")
	router.HandleFunc("/api/events", dashboardServer.handleGetEvents).Methods("GET")
	router.HandleFunc("/api/logs", dashboardServer.handleGetLogs).Methods("GET")
	router.HandleFunc("/api/users", dashboardServer.handleGetUsers).Methods("GET")
	router.HandleFunc("/api/users", dashboardServer.handleCreateUser).Methods("POST")

	// Main dashboard route
	router.HandleFunc("/", dashboardServer.handleDashboard).Methods("GET")

	log.Printf("%s", "Dashboard starting on http://localhost:"+dashboardPort)

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
