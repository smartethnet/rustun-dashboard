package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/smartethnet/rustun-dashboard/internal/ipadm"
	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/repository"
)

type RouteService struct {
	repo      repository.RouteRepository
	ipManager *ipadm.IPAdmManager
}

// NewRouteService creates a new route service with the given repository and IP manager
func NewRouteService(repo repository.RouteRepository, ipManager *ipadm.IPAdmManager) *RouteService {
	return &RouteService{
		repo:      repo,
		ipManager: ipManager,
	}
}

// GetAllClusters returns all unique clusters with client counts
func (s *RouteService) GetAllClusters() ([]model.Cluster, error) {
	clusterMap, err := s.repo.GetAllClusters()
	if err != nil {
		return nil, err
	}

	clusters := make([]model.Cluster, 0, len(clusterMap))
	for name, count := range clusterMap {
		clusters = append(clusters, model.Cluster{
			Name:        name,
			ClientCount: count,
		})
	}

	return clusters, nil
}

// GetCluster returns a specific cluster with its clients
func (s *RouteService) GetCluster(clusterName string) (*model.Cluster, []model.Client, error) {
	clients, err := s.repo.GetByCluster(clusterName)
	if err != nil {
		return nil, nil, err
	}

	if len(clients) == 0 {
		return nil, nil, fmt.Errorf("cluster not found")
	}

	cluster := &model.Cluster{
		Name:        clusterName,
		ClientCount: len(clients),
	}

	return cluster, clients, nil
}

// DeleteCluster removes all clients in a cluster
func (s *RouteService) DeleteCluster(clusterName string) error {
	return s.repo.DeleteCluster(clusterName)
}

// GetAllClients returns all clients
func (s *RouteService) GetAllClients() ([]model.Client, error) {
	return s.repo.GetAll()
}

// GetClientsByCluster returns clients in a specific cluster
func (s *RouteService) GetClientsByCluster(clusterName string) ([]model.Client, error) {
	return s.repo.GetByCluster(clusterName)
}

// GetClient returns a specific client by cluster and identity
func (s *RouteService) GetClient(clusterName, identity string) (*model.Client, error) {
	return s.repo.GetByClusterAndIdentity(clusterName, identity)
}

// CreateClient adds a new client with auto-generated identity and IP
func (s *RouteService) CreateClient(client model.Client) (*model.Client, error) {
	// Generate UUID as identity
	client.Identity = uuid.New().String()

	// Allocate IP address with network config
	allocated, err := s.ipManager.AllocateIP(client.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate IP: %w", err)
	}
	client.PrivateIP = allocated.IP
	client.Gateway = allocated.Gateway
	client.Mask = allocated.Mask

	if err := s.repo.Create(client); err != nil {
		// Release IP on failure
		s.ipManager.ReleaseIP(client.Cluster, allocated.IP)
		return nil, err
	}

	return &client, nil
}

// UpdateClient updates an existing client
func (s *RouteService) UpdateClient(clusterName, identity string, updatedClient model.Client) error {
	return s.repo.Update(clusterName, identity, updatedClient)
}

// DeleteClient removes a client and releases its IP
func (s *RouteService) DeleteClient(clusterName, identity string) error {
	// Get client to retrieve its IP
	client, err := s.repo.GetByClusterAndIdentity(clusterName, identity)
	if err != nil {
		return err
	}

	// Delete client
	if err := s.repo.Delete(clusterName, identity); err != nil {
		return err
	}

	// Release IP address
	s.ipManager.ReleaseIP(clusterName, client.PrivateIP)

	return nil
}
