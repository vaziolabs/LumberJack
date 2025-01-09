package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vaziolabs/lumberjack/internal/core"
	"github.com/vaziolabs/lumberjack/types"

	"github.com/golang-jwt/jwt"
)

// HTTP handler for assigning a user
func (server *Server) handleAssignUser(w http.ResponseWriter, r *http.Request) {
	server.logger.Enter("AssignUser")
	defer server.logger.Exit("AssignUser")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		server.logger.Failure("User ID required")
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

	node, err := server.forest.GetNode(request.Path)
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
	if err := server.writeChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for starting time tracking
func (server *Server) handleStartTimeTracking(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	node.StartTimeTracking(userID)

	// Write changes to file
	if err := server.writeChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for stopping time tracking
func (server *Server) handleStopTimeTracking(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	node.StopTimeTracking(userID)

	summary := node.GetTimeTrackingSummary(userID)
	json.NewEncoder(w).Encode(summary)

	// Write changes to file
	if err := server.writeChangesToFile(node, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}
}

// HTTP handler for getting time tracking summary
func (server *Server) handleGetTimeTracking(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.forest.GetNode(request.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	summary := node.GetTimeTrackingSummary(userID)
	json.NewEncoder(w).Encode(summary)
}

// HTTP handler for starting an event
func (server *Server) handleStartEvent(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.getNodeFromPath(request.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Path error: %v", err), http.StatusNotFound)
		return
	}

	if err := node.StartEvent(request.EventID, nil, nil, request.Metadata); err != nil {
		http.Error(w, fmt.Sprintf("Start event error: %v", err), http.StatusInternalServerError)
		return
	}

	// Save state after event creation
	if err := server.writeChangesToFile(server.forest, "state_file.dat"); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for ending an event
func (server *Server) handleEndEvent(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.getNodeFromPath(request.Path)
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
func (server *Server) handleAppendToEvent(w http.ResponseWriter, r *http.Request) {
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
	node, err := server.getNodeFromPath(request.Path)
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

	if err := server.writeChangesToFile(server.forest, "state_file.dat"); err != nil {
		log.Printf("Failed to save state: %v", err)
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for getting event entries
func (server *Server) handleGetEventEntries(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path    string `json:"path"`
		EventID string `json:"event_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := server.getNodeFromPath(request.Path)
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
func (server *Server) handleGetForest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server.forest)
}

// HTTP handler for getting users
func (server *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server.forest.Users)
}

// HTTP handler for creating a user
func (server *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	server.logger.Enter("CreateUser")
	defer server.logger.Exit("CreateUser")

	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.logger.Failure("Failed to decode request: %v", err)
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
		server.logger.Failure("Failed to set password: %v", err)
		http.Error(w, "Failed to set password", http.StatusInternalServerError)
		return
	}

	// Add user to the root node
	if err := server.forest.AssignUser(user, core.ReadPermission); err != nil {
		server.logger.Failure("Failed to assign user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save state
	if err := server.writeChangesToFile(server.forest, "state_file.dat"); err != nil {
		server.logger.Failure("Failed to save state: %v", err)
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	server.logger.Success("User created successfully")
}

func (server *Server) handlePlanEvent(w http.ResponseWriter, r *http.Request) {
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

	node, err := server.getNodeFromPath(request.Path)
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
func (server *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	server.logger.Enter("handleLogin")
	defer server.logger.Exit("handleLogin")

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		server.logger.Failure("Invalid request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	server.logger.Info("Attempting login for user: %s", credentials.Username)
	server.logger.Info("Number of users in system: %d", len(server.forest.Users))

	// Get pointer to user to avoid copying
	var foundUser *core.User
	for i := range server.forest.Users {
		if server.forest.Users[i].Username == credentials.Username {
			foundUser = &server.forest.Users[i]
			break
		}
	}

	if foundUser == nil {
		server.logger.Failure("User not found: %s", credentials.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !foundUser.VerifyPassword(credentials.Password) {
		server.logger.Failure("Invalid password for user: %s", credentials.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  foundUser.ID,
		"username": foundUser.Username,
		"exp":      time.Now().Add(server.jwtConfig.ExpiresIn).Unix(),
	})

	tokenString, err := token.SignedString(server.jwtConfig.SecretKey)
	if err != nil {
		server.logger.Failure("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    foundUser.ID,
		"token": tokenString,
	})
	server.logger.Success("Login successful for user %s", foundUser.Username)
}

func (server *Server) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	user, err := server.forest.GetUserProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username":     user.Username,
		"email":        user.Email,
		"organization": user.Organization,
		"phone":        user.Phone,
		"permissions":  user.Permissions,
	})
}

// UpdateSettings updates server configuration parameters
func (server *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings types.ServerConfig
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

	// Update server configuration
	if err := server.UpdateSettings(userID, settings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HTTP handler for getting a specific tree
func (server *Server) handleGetTree(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	node, err := server.forest.GetNode(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

// HTTP handler for getting server settings
func (server *Server) handleGetServerSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	// Check if user has admin permission on root node
	if !server.forest.CheckPermission(userID, core.AdminPermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	// Return safe subset of server settings
	settings := map[string]interface{}{
		"organization":  server.config.Organization,
		"server_port":   server.config.ServerPort,
		"dashboard_url": server.config.DashboardURL,
		"phone":         server.config.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// HTTP handler for updating server settings
func (server *Server) handleUpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	// Check if user has admin permission on root node
	if !server.forest.CheckPermission(userID, core.AdminPermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var settings types.ServerConfig
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update only safe settings
	server.config.Organization = settings.Organization
	server.config.DashboardURL = settings.DashboardURL
	server.config.Phone = settings.Phone

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
