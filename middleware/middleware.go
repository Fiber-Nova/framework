package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger returns a logging middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Log request details
		duration := time.Since(start)
		c.Append("X-Response-Time", duration.String())
		
		return err
	}
}

// CORS returns a CORS middleware
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		
		return c.Next()
	}
}

// Auth returns an authentication middleware
func Auth(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		
		// Simple token validation (in production, use proper JWT validation)
		if token != "Bearer "+key {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}
		
		return c.Next()
	}
}

// RateLimiter returns a simple rate limiting middleware
func RateLimiter(max int, window time.Duration) fiber.Handler {
	// This is a simplified version - in production use a proper rate limiter
	return func(c *fiber.Ctx) error {
		// Simple implementation - just pass through for now
		return c.Next()
	}
}
