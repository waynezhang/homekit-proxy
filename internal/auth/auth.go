package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type AuthManager struct {
	username string
	password string
	sessions map[string]*Session
}

type Session struct {
	Token     string
	ExpiresAt time.Time
}

func NewAuthManager() *AuthManager {
	username := os.Getenv("HOMEKIT_PROXY_USER")
	password := os.Getenv("HOMEKIT_PROXY_PASSWORD")
	
	if username == "" || password == "" {
		slog.Warn("[Auth] Environment variables HOMEKIT_PROXY_USER and HOMEKIT_PROXY_PASSWORD not set - authentication disabled")
		return nil
	}
	
	return &AuthManager{
		username: username,
		password: password,
		sessions: make(map[string]*Session),
	}
}

func (am *AuthManager) IsEnabled() bool {
	return am != nil
}

func (am *AuthManager) Login(username, password string) (string, error) {
	if !am.IsEnabled() {
		return "", fmt.Errorf("authentication not enabled")
	}
	
	if subtle.ConstantTimeCompare([]byte(am.username), []byte(username)) != 1 ||
		subtle.ConstantTimeCompare([]byte(am.password), []byte(password)) != 1 {
		return "", fmt.Errorf("invalid credentials")
	}
	
	// Generate session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	am.sessions[token] = &Session{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
	}
	
	// Clean up expired sessions
	am.cleanupExpiredSessions()
	
	slog.Info("[Auth] User logged in", "username", username)
	return token, nil
}

func (am *AuthManager) ValidateSession(token string) bool {
	if !am.IsEnabled() {
		return true // Authentication disabled
	}
	
	if token == "" {
		return false
	}
	
	session, exists := am.sessions[token]
	if !exists {
		return false
	}
	
	if time.Now().After(session.ExpiresAt) {
		delete(am.sessions, token)
		return false
	}
	
	return true
}

func (am *AuthManager) Logout(token string) {
	if !am.IsEnabled() {
		return
	}
	
	delete(am.sessions, token)
	slog.Info("[Auth] User logged out")
}

func (am *AuthManager) cleanupExpiredSessions() {
	now := time.Now()
	for token, session := range am.sessions {
		if now.After(session.ExpiresAt) {
			delete(am.sessions, token)
		}
	}
}

func (am *AuthManager) RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !am.IsEnabled() {
			handler(w, r)
			return
		}
		
		// Check session cookie
		cookie, err := r.Cookie("homekit_proxy_session")
		if err != nil {
			am.sendUnauthorized(w)
			return
		}
		
		if !am.ValidateSession(cookie.Value) {
			am.sendUnauthorized(w)
			return
		}
		
		handler(w, r)
	}
}

func (am *AuthManager) RequireAuthForAPI(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !am.IsEnabled() {
			handler(w, r)
			return
		}
		
		// Check session cookie first
		cookie, err := r.Cookie("homekit_proxy_session")
		if err == nil && am.ValidateSession(cookie.Value) {
			handler(w, r)
			return
		}
		
		// Check Authorization header as fallback
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if am.ValidateSession(token) {
				handler(w, r)
				return
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}
}

func (am *AuthManager) sendUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
}