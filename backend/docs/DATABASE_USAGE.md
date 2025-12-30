# Database Usage Guide

Complete guide for using MySQL database with GORM in Rustun Dashboard.

## Quick Start

### 1. Install MySQL

**macOS:**
```bash
brew install mysql
brew services start mysql
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
```

**Docker:**
```bash
docker run -d \
  --name rustun-mysql \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=rustun_dashboard \
  -e MYSQL_USER=rustun \
  -e MYSQL_PASSWORD=password123 \
  -p 3306:3306 \
  mysql:8.0
```

### 2. Create Database and User

```sql
-- Connect to MySQL
mysql -u root -p

-- Create database
CREATE DATABASE rustun_dashboard CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user
CREATE USER 'rustun'@'localhost' IDENTIFIED BY 'password123';

-- Grant privileges
GRANT ALL PRIVILEGES ON rustun_dashboard.* TO 'rustun'@'localhost';
FLUSH PRIVILEGES;
```

### 3. Update Configuration

Edit `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"

auth:
  username: "admin"
  password: "admin123"

storage:
  type: "database"  # Change from "file" to "database"
  
  database:
    type: "mysql"
    host: "localhost"
    port: 3306
    username: "rustun"
    password: "password123"
    database: "rustun_dashboard"
```

### 4. Install Dependencies

```bash
cd backend
go mod download
```

### 5. Migrate Data (Optional)

If you have existing data in `routes.json`:

```bash
# Build migration tool
go build -o bin/migrate tools/migrate.go

# Run migration
./bin/migrate -config config.yaml -json routes.json
```

### 6. Start Server

```bash
go run cmd/dashboard/main.go -config config.yaml
```

The server will automatically:
- Connect to the database
- Create tables if they don't exist
- Start the API server

## Database Schema

The application automatically creates this table:

```sql
CREATE TABLE `clients` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

## Migration Tools

### Import from JSON to Database

```bash
go run tools/migrate.go -config config.yaml -json routes.json
```

Output:
```
Connected to mysql database
Database schema migrated successfully
Found 10 clients in JSON file
✅ Imported: production/prod-gateway-01
✅ Imported: production/prod-app-server-01
...
==================================================
Migration completed!
Total:     10
Imported:  10
Errors:    0
==================================================
```

### Export from Database to JSON

```bash
go run tools/export.go -config config.yaml -output backup.json
```

This creates a JSON backup of all clients.

## Using Other Databases

### PostgreSQL

**Configuration:**
```yaml
storage:
  type: "database"
  database:
    type: "postgres"
    host: "localhost"
    port: 5432
    username: "rustun"
    password: "password123"
    database: "rustun_dashboard"
```

**Setup:**
```bash
# Install PostgreSQL
brew install postgresql
brew services start postgresql

# Create database
createdb rustun_dashboard
```

### SQLite

**Configuration:**
```yaml
storage:
  type: "database"
  database:
    type: "sqlite"
    path: "./rustun.db"
```

No setup required! The database file will be created automatically.

## API Usage (Same as File Storage)

All API endpoints work exactly the same way:

```bash
# List all clusters
curl -u admin:admin123 http://localhost:8080/api/clusters

# Create a client
curl -u admin:admin123 -X POST http://localhost:8080/api/clients \
  -H "Content-Type: application/json" \
  -d '{
    "cluster": "production",
    "identity": "new-server",
    "private_ip": "10.0.1.20",
    "mask": "255.255.255.0",
    "gateway": "10.0.1.254",
    "ciders": []
  }'

# Get client
curl -u admin:admin123 http://localhost:8080/api/clients/production/new-server

# Update client
curl -u admin:admin123 -X PUT http://localhost:8080/api/clients/production/new-server \
  -H "Content-Type: application/json" \
  -d '{
    "private_ip": "10.0.1.21",
    "mask": "255.255.255.0",
    "gateway": "10.0.1.254",
    "ciders": ["192.168.50.0/24"]
  }'

# Delete client
curl -u admin:admin123 -X DELETE http://localhost:8080/api/clients/production/new-server
```

## Troubleshooting

### Connection Issues

**Error: "Access denied for user"**
```bash
# Reset MySQL password
mysql -u root -p
ALTER USER 'rustun'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
```

**Error: "Can't connect to MySQL server"**
```bash
# Check if MySQL is running
brew services list  # macOS
systemctl status mysql  # Linux

# Start MySQL
brew services start mysql  # macOS
sudo systemctl start mysql  # Linux
```

### Migration Issues

**Error: "client already exists"**

The migration tool skips existing records. This is normal.

**Error: "table already exists"**

GORM AutoMigrate is idempotent. It's safe to run multiple times.

### Performance Tips

1. **Add indexes for frequently queried fields:**
```sql
CREATE INDEX idx_clients_private_ip ON clients(private_ip);
```

2. **Monitor slow queries:**
```yaml
# In GORM config
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

3. **Connection pooling:**
```go
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Backup and Restore

### Backup

**Using export tool:**
```bash
go run tools/export.go -config config.yaml -output backup.json
```

**Using mysqldump:**
```bash
mysqldump -u rustun -p rustun_dashboard > backup.sql
```

### Restore

**From JSON:**
```bash
go run tools/migrate.go -config config.yaml -json backup.json
```

**From SQL:**
```bash
mysql -u rustun -p rustun_dashboard < backup.sql
```

## Switching Between File and Database

**To Database:**
1. Update `config.yaml`: `storage.type: "database"`
2. Run migration: `go run tools/migrate.go`
3. Restart server

**Back to File:**
1. Export data: `go run tools/export.go -output routes.json`
2. Update `config.yaml`: `storage.type: "file"`
3. Restart server

## Development vs Production

**Development (File Storage):**
```yaml
storage:
  type: "file"
  file:
    routes_file: "./routes.json"
```

**Production (Database):**
```yaml
storage:
  type: "database"
  database:
    type: "mysql"
    host: "db.production.com"
    port: 3306
    username: "rustun"
    password: "${DB_PASSWORD}"  # Use environment variable
    database: "rustun_dashboard"
```

## Security Best Practices

1. **Use environment variables for sensitive data:**
```bash
export DB_PASSWORD="secure_password"
```

2. **Limit user privileges:**
```sql
GRANT SELECT, INSERT, UPDATE, DELETE ON rustun_dashboard.* TO 'rustun'@'localhost';
```

3. **Use SSL/TLS for remote connections:**
```yaml
database:
  type: "mysql"
  host: "remote-db.com"
  port: 3306
  username: "rustun"
  password: "password"
  database: "rustun_dashboard"
  # Add SSL config in DSN
```

4. **Regular backups:**
```bash
# Daily backup cron job
0 2 * * * mysqldump -u rustun -p'password' rustun_dashboard > /backup/rustun_$(date +\%Y\%m\%d).sql
```

## Next Steps

- [API Documentation](../README.md)
- [Configuration Guide](../config.yaml)
- [Migration Guide](DATABASE_MIGRATION.md)

