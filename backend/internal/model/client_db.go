package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ClientDB represents the database model for Client (when using ORM like GORM)
// This is prepared for future database implementation
type ClientDB struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Cluster   string    `gorm:"index:idx_cluster_identity,unique;not null" json:"cluster"`
	Identity  string    `gorm:"index:idx_cluster_identity,unique;not null" json:"identity"`
	Name      string    `gorm:"" json:"name"` // Optional friendly name
	PrivateIP string    `gorm:"not null" json:"private_ip"`
	Mask      string    `gorm:"not null" json:"mask"`
	Gateway   string    `gorm:"not null" json:"gateway"`
	Ciders    JSONArray `gorm:"type:json" json:"ciders"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ClientDB) TableName() string {
	return "clients"
}

// ToClient converts ClientDB to Client
func (c *ClientDB) ToClient() Client {
	return Client{
		Cluster:   c.Cluster,
		Identity:  c.Identity,
		Name:      c.Name,
		PrivateIP: c.PrivateIP,
		Mask:      c.Mask,
		Gateway:   c.Gateway,
		Ciders:    c.Ciders,
	}
}

// FromClient converts Client to ClientDB
func (c *ClientDB) FromClient(client Client) {
	c.Cluster = client.Cluster
	c.Identity = client.Identity
	c.Name = client.Name
	c.PrivateIP = client.PrivateIP
	c.Mask = client.Mask
	c.Gateway = client.Gateway
	c.Ciders = client.Ciders
}

// JSONArray is a custom type for storing string arrays as JSON in database
type JSONArray []string

// Scan implements sql.Scanner interface
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*j = []string{}
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface
func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

