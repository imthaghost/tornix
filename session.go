package tornix

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/google/uuid"
)

const (
	defaultSessionTimeout = 30
)

type SessionInfo struct {
	Client      *http.Client
	TorPort     int
	TorInstance *tor.Tor
}

// StartNewSession starts a new session.
func (m *Manager) StartNewSession() (string, *http.Client, error) {

	log.Print("Starting new session")
	// Check if the maximum limit has been reached

	if getSyncMapLength(m.Sessions) >= m.MaxConcurrentSessions {
		return "", nil, fmt.Errorf("maximum number of concurrent sessions reached")

	}

	// Generate a unique session ID
	var uniqueID string
	for {

		uniqueID = uuid.New().String()
		if _, exists := m.Sessions.Load(uniqueID); !exists {
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

	// Store the client and Tor instance in Sessions
	m.Sessions.Swap(uniqueID, SessionInfo{Client: client,
		TorPort: socksPort})

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
	if sessionInfo, exists := m.Sessions.Load(uniqueID); exists {
		// Close the HTTP client's connections
		if transport, ok := sessionInfo.(SessionInfo).Client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}

		// Stop the Tor instance on the session's SOCKS port
		if torInstance, ok := m.TorInstances[sessionInfo.(SessionInfo).TorPort]; ok {
			torInstance.Close()                                       // Close the Tor instance
			delete(m.TorInstances, sessionInfo.(SessionInfo).TorPort) // Remove the instance from the map
		}

		m.ReleasePort(sessionInfo.(SessionInfo).TorPort) // Release the port
		m.Sessions.Delete(uniqueID)                      // Remove the session from the map
	}
}

func getSyncMapLength(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
