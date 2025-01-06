package dashboard

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

func (s *DashboardServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("dashboard/templates/index.html"))
	tmpl.Execute(w, nil)
}

func (s *DashboardServer) handleGetTree(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("GET", s.apiEndpoint+"/get_tree", nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
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

func (s *DashboardServer) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("dashboard/templates/login.html"))
	tmpl.Execute(w, nil)
}

func (s *DashboardServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.logger.Enter("Login")
	defer s.logger.Exit("Login")

	// Read the entire request body first
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Error reading request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Forward the raw request body to the API server
	req, err := http.NewRequest("POST", s.apiEndpoint+"/login", bytes.NewBuffer(body))
	if err != nil {
		s.logger.Error("Error creating request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Error making login request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Error reading response body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response status code
	w.WriteHeader(resp.StatusCode)

	// Write response body
	w.Write(respBody)

	if resp.StatusCode == http.StatusOK {
		var loginResponse struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(respBody, &loginResponse); err != nil {
			return
		}

		// Set the auth cookie on successful login
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    loginResponse.Token,
			Path:     "/",
			HttpOnly: true,
		})
	}
}

func (s *DashboardServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			s.logger.Error("No auth token found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("Authorization", "Bearer "+cookie.Value)
		next.ServeHTTP(w, r)
	})
}
