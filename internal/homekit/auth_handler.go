package homekit

import (
	"encoding/json"
	"net/http"
	"time"
)

func (m *HMManager) startAuthHandler() {
	m.server.ServeMux().HandleFunc("/login", m.handleLogin)
	m.server.ServeMux().HandleFunc("/api/login", m.handleAPILogin)
	m.server.ServeMux().HandleFunc("/api/logout", m.handleAPILogout)
	m.server.ServeMux().HandleFunc("/api/status", m.handleAuthStatus)
}

func (m *HMManager) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "web/login.html")
		return
	}
	
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed"))
}

func (m *HMManager) handleAPILogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}
	
	if !m.authManager.IsEnabled() {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "authentication not enabled"}`))
		return
	}
	
	token, err := m.authManager.Login(req.Username, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid credentials"}`))
		return
	}
	
	// Set session cookie
	cookie := &http.Cookie{
		Name:     "homekit_proxy_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	
	http.SetCookie(w, cookie)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func (m *HMManager) handleAPILogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	
	// Get session cookie
	cookie, err := r.Cookie("homekit_proxy_session")
	if err == nil {
		m.authManager.Logout(cookie.Value)
	}
	
	// Clear session cookie
	clearCookie := &http.Cookie{
		Name:     "homekit_proxy_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	
	http.SetCookie(w, clearCookie)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func (m *HMManager) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if !m.authManager.IsEnabled() {
		w.Write([]byte(`{"enabled": false, "authenticated": true}`))
		return
	}
	
	// Check session cookie
	cookie, err := r.Cookie("homekit_proxy_session")
	authenticated := err == nil && m.authManager.ValidateSession(cookie.Value)
	
	resp := map[string]interface{}{
		"enabled":       true,
		"authenticated": authenticated,
	}
	
	json.NewEncoder(w).Encode(resp)
}