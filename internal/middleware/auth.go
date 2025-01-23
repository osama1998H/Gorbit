// internal/middleware/auth.go
package middleware

import (
	"gorbit/internal/config"
	"gorbit/internal/domain"
	"gorbit/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected creates a middleware for JWT authentication
func JWTProtected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Missing authorization header",
			})
		}

		// Check token prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token format",
			})
		}

		// Parse and validate JWT
		token, err := jwt.ParseWithClaims(tokenParts[1], &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(cfg.App.JWTSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
		}

		// Validate claims
		if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
			// Set user in context
			c.Locals("user", claims.User)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Invalid token claims",
		})
	}
}

// RoleRequired creates a middleware for role-based access control
func RoleRequired(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
		}

		if !utils.Contains(user.Roles, requiredRole) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// APIKeyAuth creates middleware for API key authentication
func APIKeyAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey != cfg.App.APIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid API key",
			})
		}

		return c.Next()
	}
}
