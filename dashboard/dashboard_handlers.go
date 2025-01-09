package dashboard

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"
)

func (s *DashboardServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("dashboard/templates/index.html"))
	tmpl.Execute(w, nil)
}

func (s *DashboardServer) handleGetTree(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("GET", s.apiEndpoint+"/forest/tree", nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error connecting to API server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (s *DashboardServer) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	// Forward request to API server
	req, _ := http.NewRequest("POST", s.apiEndpoint+"/events", r.Body)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("X-User-ID", r.Header.Get("X-User-ID"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (s *DashboardServer) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// Similar to events but for system logs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{"Log functionality to be implemented"})
}

func (s *DashboardServer) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "No authentication token found", http.StatusUnauthorized)
		return
	}

	req, _ := http.NewRequest("GET", s.apiEndpoint+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+cookie.Value)
	req.Header.Set("X-User-ID", r.Header.Get("X-User-ID"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error connecting to API server: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Always set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
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

	resp, err := http.Post(s.apiEndpoint+"/users/create", "application/json",
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

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		s.logger.Error("Error decoding login credentials: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create proper JSON payload
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		s.logger.Error("Error encoding login credentials: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make request to API server
	resp, err := http.Post(s.apiEndpoint+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("Error making login request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("Login request failed with status: %v", resp.StatusCode)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Parse the API response
	var loginResponse struct {
		Status       string `json:"status"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		s.logger.Error("Invalid response from server: %v", err)
		http.Error(w, "Invalid response from server", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Setting cookies for user: %s", credentials.Username)

	// Only set session_token and refresh_token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    loginResponse.Token,
		Path:     "/",
		HttpOnly: false, // Accessible to JS for API calls
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResponse.RefreshToken,
		Path:     "/",
		HttpOnly: true, // Not accessible to JS
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   604800, // 7 days
	})

	// Return minimal response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (s *DashboardServer) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "No refresh token found", http.StatusUnauthorized)
		return
	}

	// Call API to get new session token using refresh token
	resp, err := http.Post(s.apiEndpoint+"/refresh", "application/json",
		bytes.NewBuffer([]byte(`{"refresh_token": "`+refreshToken.Value+`"}`)))
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	var newTokens struct {
		SessionToken string `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&newTokens); err != nil {
		http.Error(w, "Invalid response", http.StatusInternalServerError)
		return
	}

	// Set new session token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    newTokens.SessionToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	w.WriteHeader(http.StatusOK)
}

func (s *DashboardServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			s.logger.Error("No session token found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("Authorization", "Bearer "+cookie.Value)
		next.ServeHTTP(w, r)
	})
}

func (s *DashboardServer) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Forward request to API server
	req, _ := http.NewRequest("GET", s.apiEndpoint+"/users/profile", nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error connecting to API server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (s *DashboardServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear both cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}

func (s *DashboardServer) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("POST", s.apiEndpoint+"/settings/update", r.Body)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("X-User-ID", r.Header.Get("X-User-ID"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error connecting to API server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
