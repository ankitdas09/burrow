package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ports    map[int]bool
	cfg      Config
}

func NewSessionManager(cfg Config) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		ports:    make(map[int]bool),
		cfg:      cfg,
	}
}

func (sm *SessionManager) allocatePort() (int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for port := sm.cfg.PortRangeMin; port <= sm.cfg.PortRangeMax; port++ {
		if sm.ports[port] {
			continue
		}
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		_ = ln.Close()
		sm.ports[port] = true
		return port, nil
	}
	return 0, fmt.Errorf("no free ports in range %d-%d", sm.cfg.PortRangeMin, sm.cfg.PortRangeMax)
}

func (sm *SessionManager) freePort(port int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.ports, port)
}

func generateID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSubdomain() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
