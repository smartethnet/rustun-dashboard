package ipadm

import (
	"fmt"
	"net"
	"sync"
)

// IPConfig represents the network configuration
type IPConfig struct {
	Network string // CIDR notation, e.g., "10.12.0.0/16"
	Gateway string // Gateway IP, e.g., "10.12.0.1"
	StartIP string // Start IP for allocation, e.g., "10.12.0.10"
	Mask    string // Subnet mask, e.g., "255.255.0.0"
}

// AllocatedIP represents an allocated IP with its network configuration
type AllocatedIP struct {
	IP      string
	Gateway string
	Mask    string
}

// IPAdmManager manages IP allocation for multiple clusters
type IPAdmManager struct {
	defaultConfig IPConfig
	clusters      map[string]*ClusterIPAlloc
	mu            sync.RWMutex
}

// ClusterIPAlloc tracks IP allocation for a single cluster
type ClusterIPAlloc struct {
	config       IPConfig
	allocatedIPs map[string]bool
	mu           sync.Mutex
}

// NewIPAdmManager creates a new IP address manager with the default config
func NewIPAdmManager(defaultConfig IPConfig) *IPAdmManager {
	return &IPAdmManager{
		defaultConfig: defaultConfig,
		clusters:      make(map[string]*ClusterIPAlloc),
	}
}

// InitFromExistingClients initializes IP allocation from existing clients
func (m *IPAdmManager) InitFromExistingClients(clients []struct {
	Cluster   string
	PrivateIP string
}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Group by cluster
	clusterIPs := make(map[string][]string)
	for _, c := range clients {
		clusterIPs[c.Cluster] = append(clusterIPs[c.Cluster], c.PrivateIP)
	}

	// Initialize each cluster with default config
	for cluster, ips := range clusterIPs {
		alloc := &ClusterIPAlloc{
			config:       m.defaultConfig,
			allocatedIPs: make(map[string]bool),
		}

		// Mark all existing IPs as allocated
		for _, ip := range ips {
			alloc.allocatedIPs[ip] = true
		}

		m.clusters[cluster] = alloc
	}
}

// AllocateIP allocates a new IP for the given cluster and returns IP with network config
func (m *IPAdmManager) AllocateIP(cluster string) (*AllocatedIP, error) {
	m.mu.Lock()
	alloc, exists := m.clusters[cluster]
	if !exists {
		// Create new cluster allocation with default config
		alloc = &ClusterIPAlloc{
			config:       m.defaultConfig,
			allocatedIPs: make(map[string]bool),
		}
		m.clusters[cluster] = alloc
	}
	m.mu.Unlock()

	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	// Start from configured start IP and find first available
	_, ipNet, _ := net.ParseCIDR(alloc.config.Network)
	currentIP := net.ParseIP(alloc.config.StartIP).To4()

	for i := 0; i < 65536; i++ { // Max attempts for /16
		ipStr := currentIP.String()

		// Check if IP is in valid range
		if !ipNet.Contains(currentIP) {
			return nil, fmt.Errorf("no available IP in cluster %s", cluster)
		}

		// Skip gateway
		if ipStr == alloc.config.Gateway {
			incrementIP(currentIP)
			continue
		}

		// Check if available
		if !alloc.allocatedIPs[ipStr] {
			// Allocate this IP
			alloc.allocatedIPs[ipStr] = true
			return &AllocatedIP{
				IP:      ipStr,
				Gateway: alloc.config.Gateway,
				Mask:    alloc.config.Mask,
			}, nil
		}

		incrementIP(currentIP)
	}

	return nil, fmt.Errorf("no available IP in cluster %s", cluster)
}

// ReleaseIP releases an IP back to the pool for the given cluster
func (m *IPAdmManager) ReleaseIP(cluster, ip string) {
	m.mu.RLock()
	alloc, exists := m.clusters[cluster]
	m.mu.RUnlock()

	if !exists {
		return
	}

	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	delete(alloc.allocatedIPs, ip)
}

// incrementIP increments an IPv4 address by 1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
