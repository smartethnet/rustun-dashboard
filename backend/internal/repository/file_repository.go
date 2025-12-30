package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/smartethnet/rustun-dashboard/internal/model"
)

// FileRepository implements RouteRepository using JSON file storage
type FileRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{
		filePath: filePath,
	}
}

// loadRoutes reads and parses the routes file
func (r *FileRepository) loadRoutes() ([]model.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read routes file: %w", err)
	}

	var routes []model.Client
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse routes file: %w", err)
	}

	return routes, nil
}

// saveRoutes writes routes to the file
func (r *FileRepository) saveRoutes(routes []model.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal routes: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write routes file: %w", err)
	}

	return nil
}

// GetAll returns all clients
func (r *FileRepository) GetAll() ([]model.Client, error) {
	return r.loadRoutes()
}

// GetByClusterAndIdentity returns a specific client
func (r *FileRepository) GetByClusterAndIdentity(cluster, identity string) (*model.Client, error) {
	routes, err := r.loadRoutes()
	if err != nil {
		return nil, err
	}

	for _, client := range routes {
		if client.Cluster == cluster && client.Identity == identity {
			return &client, nil
		}
	}

	return nil, fmt.Errorf("client not found")
}

// GetByCluster returns all clients in a cluster
func (r *FileRepository) GetByCluster(cluster string) ([]model.Client, error) {
	routes, err := r.loadRoutes()
	if err != nil {
		return nil, err
	}

	clients := make([]model.Client, 0)
	for _, client := range routes {
		if client.Cluster == cluster {
			clients = append(clients, client)
		}
	}

	return clients, nil
}

// Create adds a new client
func (r *FileRepository) Create(client model.Client) error {
	routes, err := r.loadRoutes()
	if err != nil {
		return err
	}

	// Check if client already exists
	for _, c := range routes {
		if c.Cluster == client.Cluster && c.Identity == client.Identity {
			return fmt.Errorf("client already exists")
		}
	}

	// Initialize empty ciders if nil
	if client.Ciders == nil {
		client.Ciders = []string{}
	}

	routes = append(routes, client)
	return r.saveRoutes(routes)
}

// Update updates an existing client
func (r *FileRepository) Update(cluster, identity string, updatedClient model.Client) error {
	routes, err := r.loadRoutes()
	if err != nil {
		return err
	}

	found := false
	for i, client := range routes {
		if client.Cluster == cluster && client.Identity == identity {
			// Keep the original cluster and identity
			updatedClient.Cluster = cluster
			updatedClient.Identity = identity
			if updatedClient.Ciders == nil {
				updatedClient.Ciders = []string{}
			}
			routes[i] = updatedClient
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("client not found")
	}

	return r.saveRoutes(routes)
}

// Delete removes a client
func (r *FileRepository) Delete(cluster, identity string) error {
	routes, err := r.loadRoutes()
	if err != nil {
		return err
	}

	newRoutes := make([]model.Client, 0)
	found := false
	for _, client := range routes {
		if client.Cluster == cluster && client.Identity == identity {
			found = true
			continue
		}
		newRoutes = append(newRoutes, client)
	}

	if !found {
		return fmt.Errorf("client not found")
	}

	return r.saveRoutes(newRoutes)
}

// DeleteCluster removes all clients in a cluster
func (r *FileRepository) DeleteCluster(cluster string) error {
	routes, err := r.loadRoutes()
	if err != nil {
		return err
	}

	newRoutes := make([]model.Client, 0)
	found := false
	for _, client := range routes {
		if client.Cluster != cluster {
			newRoutes = append(newRoutes, client)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("cluster not found")
	}

	return r.saveRoutes(newRoutes)
}

// GetAllClusters returns all unique clusters with counts
func (r *FileRepository) GetAllClusters() (map[string]int, error) {
	routes, err := r.loadRoutes()
	if err != nil {
		return nil, err
	}

	clusterMap := make(map[string]int)
	for _, client := range routes {
		clusterMap[client.Cluster]++
	}

	return clusterMap, nil
}
