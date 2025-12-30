package service

import (
	"fmt"

	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/repository"
)

type RouteService struct {
	repo repository.RouteRepository
}

// NewRouteService creates a new route service with the given repository
func NewRouteService(repo repository.RouteRepository) *RouteService {
	return &RouteService{
		repo: repo,
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

// CreateClient adds a new client
func (s *RouteService) CreateClient(client model.Client) error {
	return s.repo.Create(client)
}

// UpdateClient updates an existing client
func (s *RouteService) UpdateClient(clusterName, identity string, updatedClient model.Client) error {
	return s.repo.Update(clusterName, identity, updatedClient)
}

// DeleteClient removes a client
func (s *RouteService) DeleteClient(clusterName, identity string) error {
	return s.repo.Delete(clusterName, identity)
}
