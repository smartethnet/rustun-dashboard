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

// Migration tool to migrate data from JSON file to database
func main() {
	configPath := flag.String("config", "./config.yaml", "Path to config file")
	jsonFile := flag.String("json", "./routes.json", "Path to JSON file to import")
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

	// Auto-migrate schema
	if err := db.AutoMigrate(&model.ClientDB{}); err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}
	log.Println("Database schema migrated successfully")

	// Read JSON file
	data, err := os.ReadFile(*jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var clients []model.Client
	if err := json.Unmarshal(data, &clients); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	log.Printf("Found %d clients in JSON file", len(clients))

	// Import data
	successCount := 0
	errorCount := 0

	for _, client := range clients {
		var dbClient model.ClientDB
		dbClient.FromClient(client)

		// Check if already exists
		var count int64
		db.Model(&model.ClientDB{}).
			Where("cluster = ? AND identity = ?", client.Cluster, client.Identity).
			Count(&count)

		if count > 0 {
			log.Printf("⚠️  Skipping existing: %s/%s", client.Cluster, client.Identity)
			continue
		}

		if err := db.Create(&dbClient).Error; err != nil {
			log.Printf("❌ Failed to import %s/%s: %v", client.Cluster, client.Identity, err)
			errorCount++
		} else {
			log.Printf("✅ Imported: %s/%s", client.Cluster, client.Identity)
			successCount++
		}
	}

	// Summary
	fmt.Println("\n" + "="*50)
	fmt.Printf("Migration completed!\n")
	fmt.Printf("Total:     %d\n", len(clients))
	fmt.Printf("Imported:  %d\n", successCount)
	fmt.Printf("Errors:    %d\n", errorCount)
	fmt.Println("=" * 50)
}
