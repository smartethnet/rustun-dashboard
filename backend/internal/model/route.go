package model

// Client represents a VPN client configuration
type Client struct {
	Cluster   string   `json:"cluster" binding:"required"`
	Identity  string   `json:"identity" binding:"required"`
	PrivateIP string   `json:"private_ip" binding:"required"`
	Mask      string   `json:"mask" binding:"required"`
	Gateway   string   `json:"gateway" binding:"required"`
	Ciders    []string `json:"ciders"`
}

// Cluster represents a group of clients
type Cluster struct {
	Name        string `json:"name"`
	ClientCount int    `json:"client_count"`
}

// RouteConfig represents the complete routes.json structure
type RouteConfig []Client

