package shared

import "sync"

// Tunnel represents a tunnel configuration.
type Tunnel struct {
	ID       int
	Name     string
	LocalURL string
	RelayURL string
}

// InMemoryStore provides an in-memory storage for tunnels.
type InMemoryStore struct {
	mu      sync.Mutex
	tunnels map[int]Tunnel
	nextID  int
}

// NewInMemoryStore creates a new in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		tunnels: make(map[int]Tunnel),
		nextID:  1,
	}
}

// AddTunnel adds a new tunnel to the store.
func (s *InMemoryStore) AddTunnel(name, localURL, relayURL string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	tunnel := Tunnel{
		ID:       s.nextID,
		Name:     name,
		LocalURL: localURL,
		RelayURL: relayURL,
	}
	s.tunnels[s.nextID] = tunnel
	s.nextID++
	return tunnel.ID
}

// ListTunnels retrieves all tunnels from the store.
func (s *InMemoryStore) ListTunnels() []Tunnel {
	s.mu.Lock()
	defer s.mu.Unlock()
	tunnels := make([]Tunnel, 0, len(s.tunnels))
	for _, tunnel := range s.tunnels {
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

// GetTunnel retrieves a tunnel by ID.
func (s *InMemoryStore) GetTunnel(id int) (Tunnel, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tunnel, exists := s.tunnels[id]
	return tunnel, exists
}

// DeleteTunnel removes a tunnel by ID.
func (s *InMemoryStore) DeleteTunnel(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tunnels[id]; exists {
		delete(s.tunnels, id)
		return true
	}
	return false
}
