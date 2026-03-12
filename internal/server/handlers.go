package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"burrow/internal/models"
)

func (sm *SessionManager) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTunnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ClientPort == 0 {
		http.Error(w, "client_port required", http.StatusBadRequest)
		return
	}

	subdomain := req.Named
	if subdomain == "" {
		subdomain = generateSubdomain()
	}

	if req.Named != "" {
		sm.mu.RLock()
		for _, s := range sm.sessions {
			if s.Subdomain == req.Named {
				sm.mu.RUnlock()
				http.Error(w, fmt.Sprintf("named tunnel %q already active", req.Named), http.StatusConflict)
				return
			}
		}
		sm.mu.RUnlock()
	}

	serverPort, err := sm.allocatePort()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	id := generateID()
	session := &Session{
		Session: models.Session{
			ID:         id,
			Subdomain:  subdomain,
			ServerPort: serverPort,
			ClientPort: req.ClientPort,
			PublicURL:  fmt.Sprintf("https://%s.%s", subdomain, sm.cfg.BaseDomain),
			CreatedAt:  time.Now(),
		},
	}

	if err := sm.registerCaddyRoute(subdomain, serverPort); err != nil {
		sm.freePort(serverPort)
		http.Error(w, fmt.Sprintf("caddy: %v", err), http.StatusInternalServerError)
		return
	}

	sm.mu.Lock()
	sm.sessions[id] = session
	sm.mu.Unlock()

	log.Printf("[%s] created tunnel %s → localhost:%d (client port %d)", id, session.PublicURL, serverPort, req.ClientPort)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&session.Session)
}

func (sm *SessionManager) DeleteTunnel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sm.mu.Lock()
	session, ok := sm.sessions[id]
	if !ok {
		sm.mu.Unlock()
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}
	delete(sm.sessions, id)
	sm.mu.Unlock()

	_ = sm.deregisterCaddyRoute(session.Subdomain)
	sm.freePort(session.ServerPort)
	log.Printf("[%s] removed tunnel %s", id, session.PublicURL)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted", "id": id})
}

func (sm *SessionManager) ListTunnels(w http.ResponseWriter, r *http.Request) {
	sm.mu.RLock()
	list := make([]*models.Session, 0, len(sm.sessions))
	for _, s := range sm.sessions {
		list = append(list, &s.Session)
	}
	sm.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (sm *SessionManager) GetTunnel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sm.mu.RLock()
	session, ok := sm.sessions[id]
	sm.mu.RUnlock()

	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&session.Session)
}

func (sm *SessionManager) DeleteAllTunnels(w http.ResponseWriter, r *http.Request) {
	sm.mu.Lock()
	sessions := make([]*Session, 0, len(sm.sessions))
	for _, s := range sm.sessions {
		sessions = append(sessions, s)
	}
	sm.sessions = make(map[string]*Session)
	sm.mu.Unlock()

	for _, session := range sessions {
		_ = sm.deregisterCaddyRoute(session.Subdomain)
		sm.freePort(session.ServerPort)
		log.Printf("[%s] cleaned up tunnel %s", session.ID, session.PublicURL)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "cleaned",
		"removed": len(sessions),
	})
}
