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
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Post(s.apiEndpoint+"/login", "application/json",
		bytes.NewBuffer([]byte(`{"username":"`+credentials.Username+`","password":"`+credentials.Password+`"}`)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Invalid credentials", resp.StatusCode)
		return
	}

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenResp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	w.WriteHeader(http.StatusOK)
}

func (s *DashboardServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("Authorization", "Bearer "+cookie.Value)
		next.ServeHTTP(w, r)
	})
}
