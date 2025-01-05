package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
		Forest: core.NewForest("root"),
		JWTConfig: internal.JWTConfig{
			SecretKey: []byte("your-secret-key"),
			ExpiresIn: 24 * time.Hour,
		},
		Logger: internal.NewLogger(),
	}
}

// Start the HTTP server
func StartServer(app *internal.App) *http.Server {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", app.HandleLogin).Methods("POST")
	router.HandleFunc("/create_user", app.HandleCreateUser).Methods("POST")

	// Protected routes
	router.HandleFunc("/end_event", app.AuthMiddleware(app.HandleEndEvent)).Methods("POST")
	router.HandleFunc("/append_event", app.AuthMiddleware(app.HandleAppendToEvent)).Methods("POST")
	router.HandleFunc("/start_event", app.AuthMiddleware(app.HandleStartEvent)).Methods("POST")
	router.HandleFunc("/get_event_entries", app.AuthMiddleware(app.HandleGetEventEntries)).Methods("POST")
	router.HandleFunc("/plan_event", app.AuthMiddleware(app.HandlePlanEvent)).Methods("POST")
	router.HandleFunc("/assign_user", app.AuthMiddleware(app.HandleAssignUser)).Methods("POST")
	router.HandleFunc("/start_time_tracking", app.AuthMiddleware(app.HandleStartTimeTracking)).Methods("POST")
	router.HandleFunc("/stop_time_tracking", app.AuthMiddleware(app.HandleStopTimeTracking)).Methods("POST")
	router.HandleFunc("/get_time_tracking", app.AuthMiddleware(app.HandleGetTimeTracking)).Methods("GET")
	router.HandleFunc("/get_tree", app.AuthMiddleware(app.HandleGetTree)).Methods("GET")
	router.HandleFunc("/get_users", app.AuthMiddleware(app.HandleGetUsers)).Methods("GET")

	return &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
}

func main() {
	app := NewApp()
	apiServer := StartServer(app)

	// Start the dashboard server in a goroutine
	dashboardServer := dashboard.NewDashboardServer("http://localhost:8080")
	go dashboardServer.Start()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	app.Logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		app.Logger.Failure("API server forced to shutdown: %v", err)
	}

	app.Logger.Success("Servers exited properly")
	os.Exit(0)
}
