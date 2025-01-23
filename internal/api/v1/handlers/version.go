package handlers

import (
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
)

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandomNumber(c *fiber.Ctx) error {
	// Generate random number between 1-1000
	min := 1
	max := 1000
	randomNumber := rand.Intn(max-min+1) + min

	return c.JSON(fiber.Map{
		"number":  randomNumber,
		"version": "v1",
	})
}
