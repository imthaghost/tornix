package tornix

import (
	"sync"
)

// Manager manages the ports used by the proxy.
type Manager struct {
	Sessions     *sync.Map // Keeps track of active sessions
	UsedPorts    *sync.Map // Keeps track of used ports
	TorInstances *sync.Map // Keeps track of active Tor instances

	MaxConcurrentSessions int // Max concurrent sessions

}

// NewManager returns a new Manager.
func NewManager(maxSessions int) *Manager {
	portsMap := &sync.Map{}
	sessionsMap := &sync.Map{}
	torInstancesMap := &sync.Map{}

	return &Manager{
		Sessions:              sessionsMap,
		TorInstances:          torInstancesMap,
		UsedPorts:             portsMap,
		MaxConcurrentSessions: maxSessions,
	}
}

// GetSyncMapLength returns the length of a sync.Map.
func GetSyncMapLength(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
