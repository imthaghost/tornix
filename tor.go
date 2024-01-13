package tornix

import (
	"context"
	"fmt"
	"github.com/cretz/bine/tor"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

// TorProxyFunc returns a proxy function for the Tor instance running on the host machine.
func (m *Manager) TorProxyFunc(id string) (*url.URL, int, error) {

	log.Println("TorProxyFunc called")
	// Select a random Tor instance's SOCKS port
	_, socksPort, err := m.SelectRandomTorInstance()
	if err != nil {
		return nil, 0, fmt.Errorf("error selecting Tor instance: %v", err)
	}

	// Create the proxy URL using the SOCKS port and unique credentials
	proxyURLString := fmt.Sprintf("socks5://%s:x@127.0.0.1:%d", id, socksPort)
	proxyURL, err := url.Parse(proxyURLString)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing proxy URL: %w", err)
	}

	return proxyURL, socksPort, nil
}

// SelectRandomTorInstance selects a random Tor instance from the pool of active Tor instances.
func (m *Manager) SelectRandomTorInstance() (*tor.Tor, int, error) {
	log.Println("SelectRandomTorInstance called")

	// Check if there are active Tor instances
	if len(m.TorInstances) == 0 {
		// create a new Tor instance
		return m.CreateNewTorInstance()
	}

	// Existing logic to randomly select a Tor instance
	var availablePorts []int
	for port := range m.TorInstances {
		availablePorts = append(availablePorts, port)
	}

	// Randomly select an available port
	randomPort := availablePorts[rand.Intn(len(availablePorts))]

	return m.TorInstances[randomPort], randomPort, nil

}

func (m *Manager) CreateNewTorInstance() (*tor.Tor, int, error) {
	log.Println("CreateNewTorInstance called")
	// Acquire a new port for the Tor instance
	port, err := m.AcquirePort()
	if err != nil {
		return nil, 0, fmt.Errorf("unable to acquire a port: %v", err)
	}
	log.Println("Acquired port: ", port)

	// Start a new Tor instance on the acquired port
	torInstance, err := m.StartTorInstance(port)
	if err != nil {
		// Release the port if starting Tor fails
		m.ReleasePort(port)
		return nil, 0, fmt.Errorf("error starting a new Tor instance: %v", err)
	}
	log.Println("Started Tor instance on port: ", port)

	// Store the new Tor instance
	m.TorInstances[port] = torInstance
	return torInstance, port, nil
}

// StartTorInstance starts a new Tor instance on the host machine.
func (m *Manager) StartTorInstance(port int) (*tor.Tor, error) {
	dataDir := fmt.Sprintf("/tmp/tordata/%d", port) // Unique directory for each instance
	os.MkdirAll(dataDir, 0700)                      // Ensure the directory exists

	conf := &tor.StartConf{
		ExtraArgs: []string{"--SocksPort", strconv.Itoa(port)},
		DataDir:   dataDir,
	}

	t, err := tor.Start(nil, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to start Tor on port %d: %v", port, err)
	}

	//err = t.Process.Start()
	//if err != nil {
	//	return nil, fmt.Errorf("failed to start Tor on port %d: %v", port, err)
	//}
	err = t.EnableNetwork(context.Background(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to start Tor on port %d: %v", port, err)
	}
	t.DeleteDataDirOnClose = true

	return t, nil
}

func IsTorRunning(port int) bool {
	_, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 3*time.Second)
	return err == nil
}
