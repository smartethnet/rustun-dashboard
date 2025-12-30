# Database Migration Guide

This guide explains how to migrate from file-based storage to database storage.

## Architecture

The application uses the **Repository Pattern** to abstract storage operations. This allows easy switching between different storage backends without changing business logic.

### Components

1. **Repository Interface** (`internal/repository/repository.go`)
   - Defines all storage operations
   - Service layer depends on this interface

2. **File Repository** (`internal/repository/file_repository.go`)
   - Current implementation using JSON files
   - Default storage method

3. **Database Repository** (`internal/repository/database_repository.go`)
   - Placeholder for future database implementation
   - Contains TODO comments showing how to implement with GORM

## Current Implementation (File Storage)

```go
// In cmd/dashboard/main.go
repo := repository.NewFileRepository(cfg.Rustun.RoutesFile)
routeService := service.NewRouteService(repo)
```

## Migrating to Database

### Step 1: Choose a Database

Supported options:
- **MySQL** - Good for production
- **PostgreSQL** - Advanced features
- **SQLite** - Simple, file-based

### Step 2: Add Database Dependencies

Add to `go.mod`:

```bash
# For GORM (recommended ORM)
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql    # or postgres, sqlite
```

### Step 3: Implement Database Repository

Complete the implementation in `internal/repository/database_repository.go`:

```go
package repository

import (
    "fmt"
    "gorm.io/gorm"
    "github.com/smartethnet/rustun-dashboard/internal/model"
)

type DatabaseRepository struct {
    db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
    return &DatabaseRepository{db: db}
}

func (r *DatabaseRepository) GetAll() ([]model.Client, error) {
    var clients []model.ClientDB
    if err := r.db.Find(&clients).Error; err != nil {
        return nil, err
    }
    
    result := make([]model.Client, len(clients))
    for i, c := range clients {
        result[i] = c.ToClient()
    }
    return result, nil
}

// Implement other methods...
```

### Step 4: Update Configuration

Modify `config.yaml`:

```yaml
storage:
  type: "database"  # Change from "file" to "database"
  
  database:
    type: "mysql"
    host: "localhost"
    port: 3306
    username: "rustun"
    password: "your_password"
    database: "rustun_dashboard"
```

### Step 5: Update Main Application

Modify `cmd/dashboard/main.go`:

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // ... load config ...
    
    var repo repository.RouteRepository
    
    if cfg.Storage.Type == "database" {
        // Initialize database connection
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
            cfg.Storage.Database.Username,
            cfg.Storage.Database.Password,
            cfg.Storage.Database.Host,
            cfg.Storage.Database.Port,
            cfg.Storage.Database.Database,
        )
        
        db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Fatalf("Failed to connect to database: %v", err)
        }
        
        // Auto-migrate schema
        if err := db.AutoMigrate(&model.ClientDB{}); err != nil {
            log.Fatalf("Failed to migrate database: %v", err)
        }
        
        repo = repository.NewDatabaseRepository(db)
        log.Println("Using database storage")
    } else {
        // Use file storage (default)
        repo = repository.NewFileRepository(cfg.Rustun.RoutesFile)
        log.Println("Using file storage")
    }
    
    routeService := service.NewRouteService(repo)
    
    // ... rest of the code ...
}
```

### Step 6: Database Schema

The database table will be created automatically with GORM AutoMigrate:

```sql
CREATE TABLE `clients` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `cluster` varchar(255) NOT NULL,
  `identity` varchar(255) NOT NULL,
  `private_ip` varchar(255) NOT NULL,
  `mask` varchar(255) NOT NULL,
  `gateway` varchar(255) NOT NULL,
  `ciders` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cluster_identity` (`cluster`,`identity`),
  KEY `idx_clients_cluster` (`cluster`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Step 7: Data Migration

Migrate existing data from JSON file to database:

```go
// migration_tool.go
package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/smartethnet/rustun-dashboard/internal/model"
)

func main() {
    // Read JSON file
    data, err := ioutil.ReadFile("routes.json")
    if err != nil {
        log.Fatal(err)
    }
    
    var clients []model.Client
    if err := json.Unmarshal(data, &clients); err != nil {
        log.Fatal(err)
    }
    
    // Connect to database
    dsn := "user:pass@tcp(localhost:3306)/rustun_dashboard?charset=utf8mb4&parseTime=True"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Auto-migrate
    db.AutoMigrate(&model.ClientDB{})
    
    // Insert data
    for _, client := range clients {
        var dbClient model.ClientDB
        dbClient.FromClient(client)
        if err := db.Create(&dbClient).Error; err != nil {
            log.Printf("Failed to insert %s/%s: %v", client.Cluster, client.Identity, err)
        } else {
            log.Printf("Migrated %s/%s", client.Cluster, client.Identity)
        }
    }
    
    log.Println("Migration completed!")
}
```

Run migration:
```bash
go run migration_tool.go
```

## Benefits of Database Storage

1. **Concurrent Access** - Better handling of multiple simultaneous requests
2. **Indexing** - Faster queries on large datasets
3. **Transactions** - ACID guarantees for data integrity
4. **Scalability** - Better performance with many clients
5. **Queries** - Complex filtering and searching capabilities
6. **Backup** - Standard database backup tools

## Testing Both Implementations

You can keep both implementations and switch between them using configuration:

```yaml
# Development: Use file storage
storage:
  type: "file"

# Production: Use database
storage:
  type: "database"
```

## Rollback

If you need to rollback to file storage:

1. Export data from database to JSON
2. Change config back to `type: "file"`
3. Restart the application

## Best Practices

1. **Start with File Storage** - Simple and sufficient for small deployments
2. **Migrate to Database** - When you have:
   - More than 100 clients
   - High concurrent access
   - Need for complex queries
   - Multi-instance deployment

3. **Keep Repository Interface** - Never access storage directly from handlers/services
4. **Test Both Implementations** - Ensure feature parity

## Example: Complete Database Implementation

See `examples/database_repository_full.go` (to be created) for a complete working example.

