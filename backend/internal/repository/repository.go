package repository

import "github.com/smartethnet/rustun-dashboard/internal/model"

// RouteRepository defines the interface for route storage operations
type RouteRepository interface {
	// GetAll returns all clients
	GetAll() ([]model.Client, error)

	// GetByClusterAndIdentity returns a specific client
	GetByClusterAndIdentity(cluster, identity string) (*model.Client, error)

	// GetByCluster returns all clients in a cluster
	GetByCluster(cluster string) ([]model.Client, error)

	// Create adds a new client
	Create(client model.Client) error

	// Update updates an existing client
	Update(cluster, identity string, client model.Client) error

	// Delete removes a client
	Delete(cluster, identity string) error

	// DeleteCluster removes all clients in a cluster
	DeleteCluster(cluster string) error

	// GetAllClusters returns all unique clusters with counts
	GetAllClusters() (map[string]int, error)
}
