package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend-meta-data/models"
)

// ListInstruments GET /api/instruments
func ListInstruments(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ensure table exists
		if !db.Migrator().HasTable(&models.Instrument{}) {
			_ = db.AutoMigrate(&models.Instrument{})
		}
		var list []models.Instrument
		if err := db.Order("id desc").Find(&list).Error; err != nil {
			// If table missing or similar, try migrate and retry once
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "no such table") {
				_ = db.AutoMigrate(&models.Instrument{})
				if err2 := db.Order("id desc").Find(&list).Error; err2 == nil {
					return c.JSON(fiber.Map{"data": list})
				}
			}
			// fallback: return empty list to avoid 500s
			return c.JSON(fiber.Map{"data": []models.Instrument{}})
		}
		return c.JSON(fiber.Map{"data": list})
	}
}
