// @title Gorbit API
// @version 1.0
// @description API documentation for Gorbit
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorbit/internal/api"
	"gorbit/internal/api/v1/handlers"
	"gorbit/internal/cache"
	"gorbit/internal/config"
	"gorbit/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger based on environment
	var logHandler slog.Handler
	if cfg.Server.Debug {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	slog.SetDefault(slog.New(logHandler))

	// Database initialization with proper logging
	slog.Info("Initializing database connections")
	mysqlDB, err := database.InitMySQL(cfg)
	if err != nil {
		slog.Error("Failed to initialize MySQL", "error", err)
		os.Exit(1)
	}

	sqlDB, err := mysqlDB.DB()
	if err != nil {
		slog.Error("Failed to get MySQL DB instance", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Warn("Failed to close MySQL connection", "error", err)
		}
	}()

	postgresDB, err := database.InitPostgres(cfg)
	if err != nil {
		slog.Error("Failed to initialize PostgreSQL", "error", err)
		os.Exit(1)
	}

	pgDB, err := postgresDB.DB()
	if err != nil {
		slog.Error("Failed to get PostgreSQL DB instance", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := pgDB.Close(); err != nil {
			slog.Warn("Failed to close PostgreSQL connection", "error", err)
		}
	}()

	mongoDB, err := database.InitMongoDB(cfg)
	if err != nil {
		slog.Error("Failed to initialize MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDB.Disconnect(ctx); err != nil {
			slog.Warn("Failed to disconnect MongoDB", "error", err)
		}
	}()

	// Cache initialization
	slog.Debug("Initializing Redis client")
	redisClient := cache.NewRedisClient(cfg)
	if err := redisClient.Connect(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Warn("Failed to close Redis connection", "error", err)
		}
	}()

	// Create health handler
	healthHandler := handlers.NewHealthHandler(
		cfg,
		mysqlDB,
		postgresDB,
		mongoDB,
		redisClient,
	)

	// Fiber app configuration
	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name,
		ServerHeader:          fmt.Sprintf("%s v%s", cfg.App.Name, cfg.App.Version),
		DisableStartupMessage: !cfg.Server.Debug,
	})

	// Configure middleware based on environment
	if cfg.Server.Debug {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
			Output: os.Stdout,
		}))
		slog.Debug("Debug mode enabled - using verbose logging")
	} else {
		app.Use(logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
		}))
	}

	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Server.Debug,
	}))

	// Setup routes
	api.SetupRouter(app, healthHandler)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	slog.Info("Starting server",
		"address", serverAddr,
		"version", cfg.App.Version,
		"environment", environmentLabel(cfg.Server.Debug),
	)

	if err := app.Listen(serverAddr); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func environmentLabel(debug bool) string {
	if debug {
		return "development"
	}
	return "production"
}
