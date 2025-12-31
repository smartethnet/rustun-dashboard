package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smartethnet/rustun-dashboard/internal/handler"
	"github.com/smartethnet/rustun-dashboard/internal/ipadm"
	"github.com/smartethnet/rustun-dashboard/internal/middleware"
	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/repository"
	"github.com/smartethnet/rustun-dashboard/internal/service"
	"github.com/smartethnet/rustun-dashboard/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "./config.yaml", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repository based on storage type
	var repo repository.RouteRepository

	if cfg.Storage.Type == "database" {
		// Initialize database connection
		db, err := initDatabase(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// Auto-migrate schema
		if err := db.AutoMigrate(&model.ClientDB{}); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		repo = repository.NewDatabaseRepository(db)
		log.Printf("Using database storage: %s", cfg.Storage.Database.Type)
	} else {
		// Use file storage (default)
		repo = repository.NewFileRepository(cfg.Storage.File.RoutesFile)
		log.Printf("Using file storage: %s", cfg.Storage.File.RoutesFile)
	}

	// Initialize IP address manager
	ipConfig := ipadm.IPConfig{
		Network: "10.12.0.0/16",
		Gateway: "10.12.0.1",
		StartIP: "10.12.0.10",
		Mask:    "255.255.0.0",
	}
	ipManager := ipadm.NewIPAdmManager(ipConfig)
	log.Printf("Initialized IP address manager: network=%s, gateway=%s, start=%s",
		ipConfig.Network, ipConfig.Gateway, ipConfig.StartIP)

	// Initialize from existing clients
	existingClients, err := repo.GetAll()
	if err == nil && len(existingClients) > 0 {
		clientInfos := make([]struct {
			Cluster   string
			PrivateIP string
		}, len(existingClients))
		for i, client := range existingClients {
			clientInfos[i].Cluster = client.Cluster
			clientInfos[i].PrivateIP = client.PrivateIP
		}
		ipManager.InitFromExistingClients(clientInfos)
		log.Printf("Initialized IP allocations from %d existing clients", len(existingClients))
	}

	// Initialize services
	routeService := service.NewRouteService(repo, ipManager)

	// Initialize handlers
	clusterHandler := handler.NewClusterHandler(routeService)
	clientHandler := handler.NewClientHandler(routeService)

	// Setup router
	r := gin.Default()

	// Apply middleware
	r.Use(middleware.CORS())

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes with Basic Auth
	api := r.Group("/api")
	api.Use(middleware.BasicAuth(cfg.Auth.Username, cfg.Auth.Password))
	{
		// Cluster routes
		clusters := api.Group("/clusters")
		{
			clusters.GET("", clusterHandler.ListClusters)
			clusters.GET("/:name", clusterHandler.GetCluster)
			clusters.DELETE("/:name", clusterHandler.DeleteCluster)
		}

		// Client routes
		clients := api.Group("/clients")
		{
			clients.GET("", clientHandler.ListClients)
			clients.POST("", clientHandler.CreateClient)
			clients.GET("/:cluster/:identity", clientHandler.GetClient)
			clients.PUT("/:cluster/:identity", clientHandler.UpdateClient)
			clients.DELETE("/:cluster/:identity", clientHandler.DeleteClient)
		}
	}

	// Start server
	log.Printf("Starting Rustun Dashboard API server on %s", cfg.Server.Address())
	if err := r.Run(cfg.Server.Address()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDatabase initializes database connection based on config
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
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
		return nil, err
	}

	log.Printf("Successfully connected to %s database", cfg.Storage.Database.Type)
	return db, nil
}
