package tornix

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultSessionTimeout = 30
)

type SessionInfo struct {
	Client  *http.Client
	TorPort int
}

// StartNewSession starts a new session.
func (m *Manager) StartNewSession() (string, *http.Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Print("Starting new session")
	// Check if the maximum limit has been reached
	if len(m.ActiveSessions) >= m.MaxConcurrentSessions {
		return "", nil, fmt.Errorf("maximum number of concurrent sessions reached")

	}

	// Generate a unique session ID
	var uniqueID string
	for {
		uniqueID = fmt.Sprintf("user%d", rand.Int())
		if _, exists := m.ActiveSessions[uniqueID]; !exists {
			break
		}
	}
	log.Println("Unique ID: ", uniqueID)

	// Use TorProxyFunc to get a proxy URL with unique credentials
	proxyURL, socksPort, err := m.TorProxyFunc(uniqueID)
	if err != nil {
		return "", nil, fmt.Errorf("error creating proxy URL: %v", err)
	}
	log.Println("Proxy URL: ", proxyURL)

	// Create a new HTTP client with the proxy URL
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Store the client and Tor instance in ActiveSessions
	m.ActiveSessions[uniqueID] = SessionInfo{
		Client:  client,
		TorPort: socksPort,
	}

	// Start a goroutine for session timeout
	go func(uid string) {
		<-time.After(defaultSessionTimeout * time.Minute)
		m.EndSession(uid)
	}(uniqueID)

	return uniqueID, client, nil
}

// EndSession ends a session.
func (m *Manager) EndSession(uniqueID string) {
	// Assuming you have the SOCKS port associated with the session
	if sessionInfo, exists := m.ActiveSessions[uniqueID]; exists {
		// Close the HTTP client's connections
		if transport, ok := sessionInfo.Client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}

		// Stop the Tor instance on the session's SOCKS port
		if torInstance, ok := m.TorInstances[sessionInfo.TorPort]; ok {
			torInstance.Close()                         // Close the Tor instance
			delete(m.TorInstances, sessionInfo.TorPort) // Remove the instance from the map
		}

		m.ReleasePort(sessionInfo.TorPort) // Release the port
		delete(m.ActiveSessions, uniqueID) // Remove the session
	}
}
