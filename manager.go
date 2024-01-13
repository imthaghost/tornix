package tornix

import (
	"github.com/cretz/bine/tor"
	"sync"
)

// Manager manages the ports used by the proxy.
type Manager struct {
	ActiveSessions        map[string]SessionInfo // Keeps track of active sessions
	UsedPorts             map[int]bool           // Keeps track of used ports
	TorInstances          map[int]*tor.Tor       // Keeps track of active Tor instances
	MaxConcurrentSessions int                    // Max concurrent sessions

	mu sync.Mutex // Mutex for thread safety
}

// NewManager returns a new Manager.
func NewManager(maxSessions int) *Manager {
	return &Manager{
		ActiveSessions:        make(map[string]SessionInfo),
		TorInstances:          make(map[int]*tor.Tor),
		UsedPorts:             make(map[int]bool),
		MaxConcurrentSessions: maxSessions,
	}
}
