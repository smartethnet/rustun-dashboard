package repository

import (
	"fmt"

	"github.com/smartethnet/rustun-dashboard/internal/model"
	"gorm.io/gorm"
)

// DatabaseRepository implements RouteRepository using database storage with GORM
type DatabaseRepository struct {
	db *gorm.DB
}

// NewDatabaseRepository creates a new database-based repository
func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// GetAll returns all clients from database
func (r *DatabaseRepository) GetAll() ([]model.Client, error) {
	var dbClients []model.ClientDB
	if err := r.db.Find(&dbClients).Error; err != nil {
		return nil, fmt.Errorf("failed to get all clients: %w", err)
	}

	clients := make([]model.Client, len(dbClients))
	for i, dbClient := range dbClients {
		clients[i] = dbClient.ToClient()
	}

	return clients, nil
}

// GetByClusterAndIdentity returns a specific client from database
func (r *DatabaseRepository) GetByClusterAndIdentity(cluster, identity string) (*model.Client, error) {
	var dbClient model.ClientDB
	if err := r.db.Where("cluster = ? AND identity = ?", cluster, identity).First(&dbClient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	client := dbClient.ToClient()
	return &client, nil
}

// GetByCluster returns all clients in a cluster from database
func (r *DatabaseRepository) GetByCluster(cluster string) ([]model.Client, error) {
	var dbClients []model.ClientDB
	if err := r.db.Where("cluster = ?", cluster).Find(&dbClients).Error; err != nil {
		return nil, fmt.Errorf("failed to get clients by cluster: %w", err)
	}

	clients := make([]model.Client, len(dbClients))
	for i, dbClient := range dbClients {
		clients[i] = dbClient.ToClient()
	}

	return clients, nil
}

// Create adds a new client to database
func (r *DatabaseRepository) Create(client model.Client) error {
	// Check if client already exists
	var count int64
	if err := r.db.Model(&model.ClientDB{}).
		Where("cluster = ? AND identity = ?", client.Cluster, client.Identity).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing client: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("client already exists")
	}

	// Initialize empty ciders if nil
	if client.Ciders == nil {
		client.Ciders = []string{}
	}

	var dbClient model.ClientDB
	dbClient.FromClient(client)

	if err := r.db.Create(&dbClient).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

// Update updates an existing client in database
func (r *DatabaseRepository) Update(cluster, identity string, client model.Client) error {
	// Check if client exists
	var dbClient model.ClientDB
	if err := r.db.Where("cluster = ? AND identity = ?", cluster, identity).First(&dbClient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("client not found")
		}
		return fmt.Errorf("failed to find client: %w", err)
	}

	// Initialize empty ciders if nil
	if client.Ciders == nil {
		client.Ciders = []string{}
	}

	// Update fields (keep cluster and identity unchanged)
	updates := map[string]interface{}{
		"private_ip": client.PrivateIP,
		"mask":       client.Mask,
		"gateway":    client.Gateway,
		"ciders":     client.Ciders,
	}

	if err := r.db.Model(&dbClient).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}

// Delete removes a client from database
func (r *DatabaseRepository) Delete(cluster, identity string) error {
	result := r.db.Where("cluster = ? AND identity = ?", cluster, identity).Delete(&model.ClientDB{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete client: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("client not found")
	}

	return nil
}

// DeleteCluster removes all clients in a cluster from database
func (r *DatabaseRepository) DeleteCluster(cluster string) error {
	result := r.db.Where("cluster = ?", cluster).Delete(&model.ClientDB{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete cluster: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("cluster not found")
	}

	return nil
}

// GetAllClusters returns all unique clusters with counts from database
func (r *DatabaseRepository) GetAllClusters() (map[string]int, error) {
	type ClusterCount struct {
		Cluster string
		Count   int
	}

	var results []ClusterCount
	if err := r.db.Model(&model.ClientDB{}).
		Select("cluster, count(*) as count").
		Group("cluster").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get cluster counts: %w", err)
	}

	clusterMap := make(map[string]int)
	for _, result := range results {
		clusterMap[result.Cluster] = result.Count
	}

	return clusterMap, nil
}
