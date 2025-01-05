package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vaziolabs/LumberJack/internal/core"

	"github.com/golang-jwt/jwt"
)

// HTTP handler for assigning a user
func (app *App) HandleAssignUser(w http.ResponseWriter, r *http.Request) {
	app.Logger.Enter("AssignUser")
	defer app.Logger.Exit("AssignUser")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		app.Logger.Failure("User ID required")
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path       string          `json:"path"`
		AssigneeID string          `json:"assignee_id"`
		Permission core.Permission `json:"permission"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.Forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if user has admin permission
	if !node.CheckPermission(userID, core.AdminPermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	assigneeUser := core.User{ID: request.AssigneeID}
	if err := node.AssignUser(assigneeUser, request.Permission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log activity
	node.AddActivity("assign_user", map[string]interface{}{
		"assignee_id": request.AssigneeID,
		"permission":  request.Permission,
	}, userID)

	// Write changes to file
	if err := app.WriteChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for starting time tracking
func (app *App) HandleStartTimeTracking(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.Forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	node.StartTimeTracking(userID)

	// Write changes to file
	if err := app.WriteChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for stopping time tracking
func (app *App) HandleStopTimeTracking(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.Forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	node.StopTimeTracking(userID)

	summary := node.GetTimeTrackingSummary(userID)
	json.NewEncoder(w).Encode(summary)

	// Write changes to file
	if err := app.WriteChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}
}

// HTTP handler for getting time tracking summary
func (app *App) HandleGetTimeTracking(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.Forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	summary := node.GetTimeTrackingSummary(userID)
	json.NewEncoder(w).Encode(summary)
}

// HTTP handler for starting an event
func (app *App) HandleStartEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path     string                 `json:"path"`
		EventID  string                 `json:"event_id"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.getNodeFromPath(request.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Path error: %v", err), http.StatusNotFound)
		return
	}

	if err := node.StartEvent(request.EventID, nil, nil, request.Metadata); err != nil {
		http.Error(w, fmt.Sprintf("Start event error: %v", err), http.StatusInternalServerError)
		return
	}

	// Save state after event creation
	if err := app.WriteChangesToFile(app.Forest, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for ending an event
func (app *App) HandleEndEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path    string `json:"path"`
		EventID string `json:"event_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.getNodeFromPath(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := node.EndEvent(request.EventID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for appending to an event
func (app *App) HandleAppendToEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var request struct {
		Path     string                 `json:"path"`
		EventID  string                 `json:"event_id"`
		Content  string                 `json:"content"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Looking for node at path: %s", request.Path)
	node, err := app.getNodeFromPath(request.Path)
	if err != nil {
		log.Printf("Failed to get node: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Appending to event %s", request.EventID)
	if err := node.AppendToEvent(request.EventID, request.Content, request.Metadata, userID); err != nil {
		log.Printf("Failed to append to event: %v", err)
		http.Error(w, fmt.Sprintf("Failed to append to event: %v", err), http.StatusInternalServerError)
		return
	}

	if err := app.WriteChangesToFile(app.Forest, "state_file.dat"); err != nil {
		log.Printf("Failed to save state: %v", err)
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for getting event entries
func (app *App) HandleGetEventEntries(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path    string `json:"path"`
		EventID string `json:"event_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := app.getNodeFromPath(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	entries, err := node.GetEventEntries(request.EventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// HTTP handler for getting tree
func (app *App) HandleGetTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.Forest)
}

// HTTP handler for getting users
func (app *App) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.Forest.Users)
}

// HTTP handler for creating a user
func (app *App) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new user
	user := core.User{
		ID:       core.GenerateID(),
		Username: request.Username,
		Email:    request.Email,
	}

	if err := user.SetPassword(request.Password); err != nil {
		http.Error(w, "Failed to set password", http.StatusInternalServerError)
		return
	}

	// Add user to the root node
	if err := app.Forest.AssignUser(user, core.ReadPermission); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save state
	if err := app.WriteChangesToFile(app.Forest, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) HandlePlanEvent(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path      string                 `json:"path"`
		EventID   string                 `json:"event_id"`
		StartTime string                 `json:"start_time"`
		EndTime   string                 `json:"end_time"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	node, err := app.getNodeFromPath(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := node.PlanEvent(request.EventID, &startTime, &endTime, request.Metadata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Add new handlers
func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find user
	var user *core.User
	for _, u := range app.Forest.Users {
		if u.Username == credentials.Username {
			user = &u
			break
		}
	}

	if user == nil || !user.VerifyPassword(credentials.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(app.JWTConfig.ExpiresIn).Unix(),
	})

	tokenString, err := token.SignedString(app.JWTConfig.SecretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
