package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Export tool to export data from database to JSON file
func main() {
	configPath := flag.String("config", "./config.yaml", "Path to config file")
	outputFile := flag.String("output", "./routes_export.json", "Path to output JSON file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	var dialector gorm.Dialector
	switch cfg.Storage.Database.Type {
	case "mysql":
		dialector = mysql.Open(cfg.Storage.Database.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.Storage.Database.DSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.Storage.Database.DSN())
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Storage.Database.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Connected to %s database", cfg.Storage.Database.Type)

	// Read all clients from database
	var dbClients []model.ClientDB
	if err := db.Find(&dbClients).Error; err != nil {
		log.Fatalf("Failed to read clients: %v", err)
	}
	log.Printf("Found %d clients in database", len(dbClients))

	// Convert to Client model
	clients := make([]model.Client, len(dbClients))
	for i, dbClient := range dbClients {
		clients[i] = dbClient.ToClient()
	}

	// Write to JSON file
	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(*outputFile, data, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	// Summary
	fmt.Println("\n" + "="*50)
	fmt.Printf("Export completed!\n")
	fmt.Printf("Total clients: %d\n", len(clients))
	fmt.Printf("Output file:   %s\n", *outputFile)
	fmt.Println("=" * 50)
}
