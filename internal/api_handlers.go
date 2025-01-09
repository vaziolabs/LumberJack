package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vaziolabs/lumberjack/internal/core"
	"github.com/vaziolabs/lumberjack/types"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// HTTP handler for assigning a user
func (server *Server) handleAssignUser(w http.ResponseWriter, r *http.Request) {
	server.logger.Enter("AssignUser")
	defer server.logger.Exit("AssignUser")

	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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

	if err := node.StartEvent(request.EventID, userID, nil, nil, request.Metadata); err != nil {
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
	userID := r.Context().Value("user_id").(string)

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

	if !node.CheckPermission(userID, core.WritePermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := node.EndEvent(request.EventID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HTTP handler for appending to an event
func (server *Server) handleAppendToEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

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
	entry := core.Entry{
		Content:   request.Content,
		Metadata:  request.Metadata,
		UserID:    userID,
		Timestamp: time.Now(),
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	if err := node.AppendToEvent(request.EventID, userID, entry, request.Metadata); err != nil {
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
	userID := r.Context().Value("user_id").(string)

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

	if err := node.PlanEvent(request.EventID, userID, &startTime, &endTime, request.Metadata); err != nil {
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

	// Generate token pair
	tokenPair, err := server.generateTokenPair(foundUser)
	if err != nil {
		server.logger.Failure("Failed to generate tokens: %v", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_token": tokenPair.SessionToken,
		"refresh_token": tokenPair.RefreshToken,
	})
	server.logger.Success("Login successful for user %s", foundUser.Username)
}

// Add new handler for token refresh
func (server *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(request.RefreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return server.jwtConfig.SecretKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || claims.TokenType != "refresh" {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	// Generate new session token
	user := &core.User{ID: claims.UserID, Username: claims.Username}
	tokenPair, err := server.generateTokenPair(user)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"session_token": tokenPair.SessionToken,
	})
}

func (server *Server) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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
	userID := r.Context().Value("user_id").(string)

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

	// Use the UpdateSettings helper instead of direct assignment
	if err := server.UpdateSettings(userID, settings); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update settings: %v", err), http.StatusInternalServerError)
		return
	}

	// Save state after settings update
	if err := server.writeChangesToFile(server.forest, server.config.DatabasePath); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Add these structures for token management
type TokenPair struct {
	SessionToken string `json:"session_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // "session" or "refresh"
	jwt.StandardClaims
}

// handleUploadAttachment handles file uploads and creates attachments
func (server *Server) handleUploadAttachment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := r.FormValue("path")
	node, err := server.forest.GetNode(path)
	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	if !node.CheckPermission(userID, core.WritePermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	attachment := &core.Attachment{
		ID:         fmt.Sprintf("att-%d", time.Now().UnixNano()),
		Name:       header.Filename,
		Type:       header.Header.Get("Content-Type"),
		Size:       header.Size,
		UploadedBy: userID,
		UploadedAt: time.Now(),
	}

	if err := node.AddAttachment(attachment, userID); err != nil {
		http.Error(w, "Failed to add attachment to node", http.StatusInternalServerError)
		return
	}

	// Save state after attachment upload
	if err := server.writeChangesToFile(server.forest, server.config.DatabasePath); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachment)
}

// handleGetAttachment retrieves an attachment
func (server *Server) handleGetAttachment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	attachmentID := vars["id"]
	path := r.URL.Query().Get("path")

	node, err := server.forest.GetNode(path)
	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	if !node.CheckPermission(userID, core.ReadPermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	attachment, err := node.GetAttachment(attachmentID)
	if err != nil {
		http.Error(w, "Attachment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", attachment.Type)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachment.Name))
	w.Write(attachment.Data)
}

// handleAddEntryAttachment adds an attachment to a specific event entry
func (server *Server) handleAddEntryAttachment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	eventID := vars["eventId"]
	entryIndex := vars["entryIndex"]

	path := r.URL.Query().Get("path")
	node, err := server.forest.GetNode(path)
	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	if !node.CheckPermission(userID, core.WritePermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	index, err := strconv.Atoi(entryIndex)
	if err != nil {
		http.Error(w, "Invalid entry index", http.StatusBadRequest)
		return
	}

	attachment, err := core.NewAttachmentStore().Store(file, header, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to store attachment: %v", err), http.StatusInternalServerError)
		return
	}

	if err := node.AddEntryAttachment(eventID, index, attachment, userID); err != nil {
		http.Error(w, "Failed to add attachment to entry", http.StatusInternalServerError)
		return
	}

	// Save state
	if err := server.writeChangesToFile(server.forest, server.config.DatabasePath); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachment)
}

// handleDeleteAttachment deletes an attachment
func (server *Server) handleDeleteAttachment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	attachmentID := vars["id"]
	path := r.URL.Query().Get("path")

	node, err := server.forest.GetNode(path)
	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	if !node.CheckPermission(userID, core.WritePermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := node.DeleteAttachment(attachmentID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete attachment: %v", err), http.StatusInternalServerError)
		return
	}

	// Save state after deletion
	if err := server.writeChangesToFile(server.forest, server.config.DatabasePath); err != nil {
		http.Error(w, "Failed to save state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) readLogFile(path string, level string) ([]types.LogEntry, error) {
	server.logger.Info("Reading log file: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var logs []types.LogEntry
	scanner := bufio.NewScanner(file)

	// Map of symbols to log levels
	levelMap := map[string]string{
		"ℹ": "INFO",
		"✓": "SUCCESS",
		"✗": "FAILURE",
		"🔍": "DEBUG",
		"📝": "NOTICE",
		"⚠": "WARNING",
		"❌": "ERROR",
		"🔥": "CRITICAL",
		"🚨": "ALERT",
		"💀": "EMERGENCY",
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse timestamp (first 19 characters: "2025/01/09 13:02:06")
		if len(line) < 19 {
			continue
		}

		timestamp, err := time.Parse("2006/01/02 15:04:05", line[:19])
		if err != nil {
			continue
		}

		// Remove timestamp and get remainder
		remainder := strings.TrimSpace(line[19:])

		// Remove tree characters and spaces
		remainder = strings.TrimLeft(remainder, "│└┌─ ")

		var logLevel, message string

		// Check for BEGIN/END messages
		if strings.HasPrefix(remainder, "BEGIN:") {
			logLevel = "INFO"
			message = "Started: " + strings.TrimSpace(strings.TrimPrefix(remainder, "BEGIN:"))
		} else if strings.HasPrefix(remainder, "END:") {
			logLevel = "INFO"
			message = "Completed: " + strings.TrimSpace(strings.TrimPrefix(remainder, "END:"))
		} else {
			// Check for level symbols
			found := false
			for symbol, level := range levelMap {
				if strings.HasPrefix(remainder, symbol) {
					logLevel = level
					message = strings.TrimSpace(strings.TrimPrefix(remainder, symbol))
					found = true
					break
				}
			}

			if !found {
				logLevel = "INFO"
				message = remainder
			}
		}

		// Filter by level if specified
		if level != "" && !strings.EqualFold(level, logLevel) {
			continue
		}

		entry := types.LogEntry{
			Timestamp: timestamp,
			Level:     logLevel,
			Message:   message,
		}

		logs = append(logs, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %v", err)
	}

	server.logger.Info("Found %d log entries", len(logs))
	return logs, nil
}

func (server *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	server.logger.Enter("handleGetLogs")
	defer server.logger.Exit("handleGetLogs")

	userID := r.Context().Value("user_id").(string)

	if !server.forest.CheckPermission(userID, core.AdminPermission) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	server.logger.Info("Getting logs for user %s", userID)

	level := r.URL.Query().Get("level")

	// Use the process ID from server config
	logPath := filepath.Join(server.config.LogDirectory, fmt.Sprintf("%s.log", server.config.ProcessInfo.ID))

	logs, err := server.readLogFile(logPath, level)
	if err != nil {
		server.logger.Failure("Failed to read logs: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
