package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"dashboard"

	"github.com/vaziolabs/LumberJack/internal"
	"github.com/vaziolabs/LumberJack/internal/core"

	"github.com/gorilla/mux"
)

// Initialize the application
func NewApp() *internal.App {
	return &internal.App{
		forest: core.NewForest("root"),
		jwtConfig: internal.JWTConfig{
			SecretKey: []byte("your-secret-key"),
			ExpiresIn: 24 * time.Hour,
		},
		logger: internal.NewLogger(),
	}
}

// Start the HTTP server
func (app *App) StartServer() *http.Server {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", app.handleLogin).Methods("POST")
	router.HandleFunc("/create_user", app.handleCreateUser).Methods("POST")

	// Protected routes
	router.HandleFunc("/end_event", app.authMiddleware(app.handleEndEvent)).Methods("POST")
	router.HandleFunc("/append_event", app.authMiddleware(app.handleAppendToEvent)).Methods("POST")
	router.HandleFunc("/start_event", app.authMiddleware(app.handleStartEvent)).Methods("POST")
	router.HandleFunc("/get_event_entries", app.authMiddleware(app.handleGetEventEntries)).Methods("POST")
	router.HandleFunc("/plan_event", app.authMiddleware(app.handlePlanEvent)).Methods("POST")
	router.HandleFunc("/assign_user", app.authMiddleware(app.handleAssignUser)).Methods("POST")
	router.HandleFunc("/start_time_tracking", app.authMiddleware(app.handleStartTimeTracking)).Methods("POST")
	router.HandleFunc("/stop_time_tracking", app.authMiddleware(app.handleStopTimeTracking)).Methods("POST")
	router.HandleFunc("/get_time_tracking", app.authMiddleware(app.handleGetTimeTracking)).Methods("GET")
	router.HandleFunc("/get_tree", app.authMiddleware(app.handleGetTree)).Methods("GET")
	router.HandleFunc("/get_users", app.authMiddleware(app.handleGetUsers)).Methods("GET")

	return &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
}

func main() {
	app := NewApp()
	apiServer := app.StartServer()

	// Start the dashboard server in a goroutine
	dashboardServer := dashboard.NewDashboardServer("http://localhost:8080")
	go dashboardServer.Start()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	app.logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		app.logger.Failure("API server forced to shutdown: %v", err)
	}

	app.logger.Success("Servers exited properly")
	os.Exit(0)
}

// getNodeFromPath traverses the forest to find a node by its path
func (app *App) getNodeFromPath(path string) (*forestree.Node, error) {
	if path == "" {
		return app.forest, nil
	}

	parts := strings.Split(path, "/")
	current := app.forest

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
