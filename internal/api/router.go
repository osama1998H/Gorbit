package api

import (
	v1 "gorbit/internal/api/v1"

	"github.com/gofiber/fiber/v2"

	// "gorbit/internal/config"
	"gorbit/internal/api/v1/handlers"
)

func SetupRouter(app *fiber.App, healthHandler *handlers.HealthHandler) {
	apiGroup := app.Group("/api")
	v1.RegisterRoutes(apiGroup, healthHandler)
}
