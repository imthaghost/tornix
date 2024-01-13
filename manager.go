package tornix

import (
	"sync"

	"github.com/cretz/bine/tor"
)

// Manager manages the ports used by the proxy.
type Manager struct {
	Sessions              *sync.Map        // Keeps track of active sessions
	UsedPorts             *sync.Map        // Keeps track of used ports
	TorInstances          map[int]*tor.Tor // Keeps track of active Tor instances
	MaxConcurrentSessions int              // Max concurrent sessions

}

// NewManager returns a new Manager.
func NewManager(maxSessions int) *Manager {
	portsMap := &sync.Map{}

	sessionsMap := &sync.Map{}
	return &Manager{
		Sessions:              sessionsMap,
		TorInstances:          make(map[int]*tor.Tor),
		UsedPorts:             portsMap,
		MaxConcurrentSessions: maxSessions,
	}
}
