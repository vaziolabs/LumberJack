package dashboard

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DashboardServer struct {
	apiEndpoint string
}

type TreeNode struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Type     string               `json:"type"`
	Children map[string]*TreeNode `json:"children"`
	Events   map[string]EventData `json:"events"`
	Entries  []EntryData          `json:"entries"`
}

type EventData struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Status    string     `json:"status"`
	Category  string     `json:"category"`
}

type EntryData struct {
	Content   interface{} `json:"content"`
	UserID    string      `json:"user_id"`
	Timestamp time.Time   `json:"timestamp"`
}

type UserData struct {
	ID          string       `json:"id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Permissions []Permission `json:"permissions"`
}

type Permission int

const (
	ReadPermission Permission = iota
	WritePermission
	AdminPermission
)

func NewDashboardServer(apiEndpoint string) *DashboardServer {
	return &DashboardServer{
		apiEndpoint: apiEndpoint,
	}
}

func (s *DashboardServer) Start() *http.Server {
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("dashboard/static"))))

	// API routes
	router.HandleFunc("/api/tree", s.handleGetTree).Methods("GET")
	router.HandleFunc("/api/events", s.handleGetEvents).Methods("GET")
	router.HandleFunc("/api/logs", s.handleGetLogs).Methods("GET")
	router.HandleFunc("/api/users", s.handleGetUsers).Methods("GET")
	router.HandleFunc("/api/users", s.handleCreateUser).Methods("POST")

	// Main dashboard route
	router.HandleFunc("/", s.handleDashboard).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         ":8081",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Dashboard starting on http://localhost:8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Dashboard server error: %v\n", err)
		}
	}()

	return server
}

func (s *DashboardServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("dashboard/templates/index.html"))
	tmpl.Execute(w, nil)
}

func (s *DashboardServer) handleGetTree(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(s.apiEndpoint + "/get_tree")
	if err != nil {
		http.Error(w, "Error connecting to API server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if response is an error message
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response directly
	w.Write(body)
}

func (s *DashboardServer) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(s.apiEndpoint + "/get_event_entries/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var events []EntryData
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (s *DashboardServer) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// Similar to events but for system logs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{"Log functionality to be implemented"})
}

func (s *DashboardServer) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(s.apiEndpoint + "/get_users")
	if err != nil {
		http.Error(w, "Error connecting to API server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if response is an error message
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Write the response directly
	w.Write(body)
}

func (s *DashboardServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Post(s.apiEndpoint+"/create_user", "application/json",
		bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
