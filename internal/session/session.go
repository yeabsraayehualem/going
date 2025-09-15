package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"

	"going/internal/config"
)

// Session represents a user session
type Session struct {
	ID        string
	Values    map[string]interface{}
	ExpiresAt time.Time
}

// Manager handles session creation and management
type Manager struct {
	config     *config.Config
	sessions   map[string]*Session
	mu         sync.RWMutex
	expiration time.Duration
}

// NewManager creates a new session manager
func NewManager(cfg *config.Config) *Manager {
	expiration := time.Duration(cfg.Session.Lifetime) * time.Minute
	return &Manager{
		config:     cfg,
		sessions:   make(map[string]*Session),
		expiration: expiration,
	}
}

// CreateSession creates a new session
func (m *Manager) CreateSession() *Session {
	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		Values:    make(map[string]interface{}),
		ExpiresAt: time.Now().Add(m.expiration),
	}

	m.mu.Lock()
	m.sessions[sessionID] = session
	m.mu.Unlock()

	// Start a goroutine to clean up expired sessions
	go m.cleanupExpiredSessions()

	return session
}

// GetSession retrieves a session by ID
func (m *Manager) GetSession(sessionID string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[sessionID]
	m.mu.RUnlock()

	if !exists || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session not found or expired")
	}

	// Update the expiration time on access
	session.ExpiresAt = time.Now().Add(m.expiration)
	return session, nil
}

// DeleteSession removes a session
func (m *Manager) DeleteSession(sessionID string) {
	m.mu.Lock()
	delete(m.sessions, sessionID)
	m.mu.Unlock()
}

// GetSessionFromRequest gets the session from an HTTP request
func (m *Manager) GetSessionFromRequest(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(m.config.Session.Name)
	if err != nil {
		return nil, err
	}

	return m.GetSession(cookie.Value)
}

// SetSessionCookie sets the session cookie on the response
func (m *Manager) SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.config.Session.Name,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(m.expiration),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie
func (m *Manager) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.config.Session.Name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

// cleanupExpiredSessions removes expired sessions
func (m *Manager) cleanupExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		if session.ExpiresAt.Before(now) {
			delete(m.sessions, id)
		}
	}
}

// generateSessionID generates a random session ID
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.URLEncoding.EncodeToString(b)
}
