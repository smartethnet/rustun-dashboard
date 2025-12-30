# Rustun Dashboard Backend

RESTful API backend for Rustun VPN dashboard, built with Go and Gin framework.

## Features

- ✅ RESTful API design
- ✅ Basic Authentication (from config file)
- ✅ Cluster management
- ✅ Client management
- ✅ CORS support
- ✅ **Repository Pattern** - Easy to switch between file and database storage
- ✅ File-based storage (default)
- 🔄 Database storage (extensible, see [Database Migration Guide](docs/DATABASE_MIGRATION.md))

## Architecture

The application uses **Repository Pattern** for storage abstraction:

```
┌─────────────┐
│  Handlers   │  - HTTP request handling
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │  - Business logic
└──────┬──────┘
       │
┌──────▼──────────┐
│  Repository     │  - Storage interface
│  (Interface)    │
└──────┬──────────┘
       │
       ├────────────┐
       │            │
┌──────▼────┐  ┌───▼──────┐
│   File    │  │ Database │
│ Storage   │  │ Storage  │
└───────────┘  └──────────┘
```

**Benefits:**
- Easy to switch storage backends
- Better testability
- Clean separation of concerns
- Future-proof architecture

## Storage Options

### Current: File Storage (Default)

Uses JSON file (`routes.json`) for data persistence.

**Pros:**
- Simple setup
- No dependencies
- Good for small deployments (<100 clients)
- Easy to backup and version control

**Cons:**
- Not ideal for concurrent writes
- Limited query capabilities
- Can be slow with large datasets

### Future: Database Storage

Can be migrated to use database (MySQL, PostgreSQL, SQLite).

**See [Database Migration Guide](docs/DATABASE_MIGRATION.md) for details.**

**Pros:**
- Better concurrency
- Faster queries with indexing
- ACID transactions
- Scalable

## API Endpoints

### Authentication

All API endpoints (except `/health`) require Basic Authentication.

**Default Credentials:**
- Username: `admin`
- Password: `admin123`

Configure in `config.yaml`:

```yaml
auth:
  username: "admin"
  password: "admin123"
```

### Health Check

```
GET /health
```

No authentication required.

### Clusters

#### List all clusters

```
GET /api/clusters
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "name": "production",
      "client_count": 3
    },
    {
      "name": "development",
      "client_count": 2
    }
  ]
}
```

#### Get cluster details

```
GET /api/clusters/{name}
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "cluster": {
      "name": "production",
      "client_count": 2
    },
    "clients": [
      {
        "cluster": "production",
        "identity": "prod-gateway-01",
        "private_ip": "10.0.1.1",
        "mask": "255.255.255.0",
        "gateway": "10.0.1.254",
        "ciders": ["192.168.100.0/24"]
      }
    ]
  }
}
```

#### Delete cluster

```
DELETE /api/clusters/{name}
```

Deletes the cluster and all its clients.

### Clients

#### List all clients

```
GET /api/clients
GET /api/clients?cluster=production  # Filter by cluster
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "cluster": "production",
      "identity": "prod-gateway-01",
      "private_ip": "10.0.1.1",
      "mask": "255.255.255.0",
      "gateway": "10.0.1.254",
      "ciders": ["192.168.100.0/24"]
    }
  ]
}
```

#### Get client

```
GET /api/clients/{cluster}/{identity}
```

#### Create client

```
POST /api/clients
Content-Type: application/json

{
  "cluster": "production",
  "identity": "new-client-01",
  "private_ip": "10.0.1.10",
  "mask": "255.255.255.0",
  "gateway": "10.0.1.254",
  "ciders": []
}
```

#### Update client

```
PUT /api/clients/{cluster}/{identity}
Content-Type: application/json

{
  "private_ip": "10.0.1.11",
  "mask": "255.255.255.0",
  "gateway": "10.0.1.254",
  "ciders": ["192.168.200.0/24"]
}
```

#### Delete client

```
DELETE /api/clients/{cluster}/{identity}
```

## Configuration

Create `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # or "release"

auth:
  username: "admin"
  password: "admin123"

storage:
  type: "file"  # or "database" (future)
  
  file:
    routes_file: "/etc/rustun/routes.json"
    routes_file_fallback: "./routes.json"
```

## Development

### Prerequisites

- Go 1.21 or higher

### Install dependencies

```bash
go mod download
```

### Run in development mode

```bash
make dev
# or
go run cmd/dashboard/main.go -config config.yaml
```

### Build

```bash
make build
```

The binary will be created at `bin/dashboard`.

### Run tests

```bash
make test
```

## Project Structure

```
backend/
├── cmd/
│   └── dashboard/              # Main application entry
│       └── main.go
├── internal/
│   ├── handler/                # HTTP request handlers
│   │   ├── cluster_handler.go
│   │   └── client_handler.go
│   ├── service/                # Business logic
│   │   └── route_service.go
│   ├── repository/             # Storage abstraction
│   │   ├── repository.go          # Interface
│   │   ├── file_repository.go     # File implementation
│   │   └── database_repository.go # Database implementation (TODO)
│   ├── model/                  # Data models
│   │   ├── route.go
│   │   ├── response.go
│   │   └── client_db.go           # Database model
│   └── middleware/             # HTTP middleware
│       ├── auth.go
│       └── cors.go
├── pkg/
│   └── config/                 # Configuration management
│       └── config.go
├── docs/
│   └── DATABASE_MIGRATION.md   # How to migrate to database
├── config.yaml                 # Configuration file
├── routes.json                 # Routes data (file storage)
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Example Usage with curl

### List clusters

```bash
curl -u admin:admin123 http://localhost:8080/api/clusters
```

### Create a client

```bash
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
```

### Update a client

```bash
curl -u admin:admin123 -X PUT http://localhost:8080/api/clients/production/new-server \
  -H "Content-Type: application/json" \
  -d '{
    "private_ip": "10.0.1.21",
    "mask": "255.255.255.0",
    "gateway": "10.0.1.254",
    "ciders": ["192.168.50.0/24"]
  }'
```

### Delete a client

```bash
curl -u admin:admin123 -X DELETE http://localhost:8080/api/clients/production/new-server
```

## Extending to Database

Want to use a database instead of files? See the [Database Migration Guide](docs/DATABASE_MIGRATION.md) for step-by-step instructions.

The application is designed to make this transition seamless!

## License

MIT
