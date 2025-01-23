package v1

import (
	"gorbit/internal/api/v1/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, healthHandler *handlers.HealthHandler) {
	// Health Check
	// router.Get("/health", healthHandler.HealthCheck)
	v1Group := router.Group("/v1")
	v1Group.Get("/health", healthHandler.HealthCheck)

	// Add other routes here
	// router.Get("/users", handlers.GetUsers)
}
