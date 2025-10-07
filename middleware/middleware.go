package middleware

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	// Console logging without timestamp
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
		DisableColors:          false,
	})
	log.SetLevel(log.InfoLevel)
}

// LoggerMiddleware logs each request using logrus
func LoggerMiddleware(c *fiber.Ctx) error {
	log.WithFields(log.Fields{
		"method": c.Method(),
		"url":    c.OriginalURL(),
		"ip":     c.IP(),
	}).Info("Traffic: Incoming request")
	return c.Next()
}
