package tornix

import (
	"errors"
	"fmt"
	"net"

	"github.com/cretz/bine/tor"
)

// ReleasePort releases a port for use by the proxy.
func (m *Manager) ReleasePort(port int) {
	// Close any active Tor instance using this port
	if torInstance, exists := m.TorInstances.Load(port); exists {
		// Assume that torInstance has a method to stop the Tor process
		// You'll need to implement this based on how you're managing Tor instances
		torInstance.(*tor.Tor).Close() // This is a placeholder. Replace with actual method to stop the Tor process

		// Remove the instance from the map after stopping it
		m.TorInstances.Delete(port)
	}

	// Mark the port as available again
	m.UsedPorts.Delete(port)

}

// IsPortAvailable returns true if the port is available for use.
func IsPortAvailable(port int) bool {
	// Check if the port is already in use
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	// close the listener when done
	ln.Close()

	return true
}

// AcquirePort acquires a port for use by the proxy.
func (m *Manager) AcquirePort() (int, error) {
	for port := 10000; port < 65535; port++ {

		_, inUse := m.UsedPorts.Load(port)
		// check if port is available and not in use
		if IsPortAvailable(port) && !inUse {
			m.UsedPorts.Store(port, true)
			return port, nil
		}

	}

	return 0, errors.New("no available ports")
}
