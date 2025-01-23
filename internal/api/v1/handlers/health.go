// @Summary Health check
// @Description Get service health status
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/v1/health [get]
// internal/api/v1/handlers/health.go
package handlers

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"gorbit/internal/cache"
	"gorbit/internal/config"

	// "gorbit/internal/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type HealthHandler struct {
	cfg         *config.Config
	mysqlDB     *gorm.DB
	postgresDB  *gorm.DB
	mongoClient *mongo.Client
	redisClient *cache.RedisClient
	startTime   time.Time
}

func NewHealthHandler(
	cfg *config.Config,
	mysqlDB *gorm.DB,
	postgresDB *gorm.DB,
	mongoClient *mongo.Client,
	redisClient *cache.RedisClient,
) *HealthHandler {
	return &HealthHandler{
		cfg:         cfg,
		mysqlDB:     mysqlDB,
		postgresDB:  postgresDB,
		mongoClient: mongoClient,
		redisClient: redisClient,
		startTime:   time.Now().UTC(),
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	System    SystemStats       `json:"system"`
	Services  map[string]string `json:"services"`
}

type SystemStats struct {
	GoVersion    string `json:"go_version"`
	Memory       string `json:"memory"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	services := make(map[string]string)
	ctx := context.Background()

	// Check MySQL
	services["mysql"] = h.checkMySQL(ctx)

	// Check PostgreSQL
	services["postgres"] = h.checkPostgreSQL(ctx)

	// Check MongoDB
	services["mongodb"] = h.checkMongoDB(ctx)

	// Check Redis
	services["redis"] = h.checkRedis(ctx)

	// Get system stats
	sysStats := h.getSystemStats()

	// Determine overall status
	overallStatus := h.determineOverallStatus(services)

	response := HealthResponse{
		Status:    overallStatus,
		Version:   h.cfg.App.Version,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(h.startTime).Truncate(time.Second).String(),
		System:    sysStats,
		Services:  services,
	}

	return c.Status(h.statusCode(overallStatus)).JSON(response)
}

func (h *HealthHandler) checkMySQL(ctx context.Context) string {
	sqlDB, err := h.mysqlDB.DB()
	if err != nil {
		return statusString(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return statusString(sqlDB.PingContext(ctx))
}

func (h *HealthHandler) checkPostgreSQL(ctx context.Context) string {
	pgDB, err := h.postgresDB.DB()
	if err != nil {
		return statusString(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return statusString(pgDB.PingContext(ctx))
}

func (h *HealthHandler) checkMongoDB(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return statusString(h.mongoClient.Ping(ctx, nil))
}

func (h *HealthHandler) checkRedis(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Access through the wrapper client
	_, err := h.redisClient.GetClient().Ping(ctx).Result()
	return statusString(err)
}

func (h *HealthHandler) getSystemStats() SystemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemStats{
		GoVersion:    runtime.Version(),
		Memory:       fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func (h *HealthHandler) determineOverallStatus(services map[string]string) string {
	for _, status := range services {
		if status != "healthy" {
			return "degraded"
		}
	}
	return "healthy"
}

func (h *HealthHandler) statusCode(overallStatus string) int {
	if overallStatus == "degraded" {
		return fiber.StatusServiceUnavailable
	}
	return fiber.StatusOK
}

func statusString(err error) string {
	if err != nil {
		return fmt.Sprintf("unhealthy: %v", err)
	}
	return "healthy"
}
