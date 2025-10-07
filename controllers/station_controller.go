package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend-meta-data/models"
)

// ListStations GET /api/station
func ListStations(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !db.Migrator().HasTable(&models.Station{}) {
			_ = db.AutoMigrate(&models.Station{})
		}
		var list []models.Station
		query := db.Order("id asc")

		if typeParam := c.Query("type"); typeParam != "" {
			query = query.Where("station_type_id = ?", typeParam)
		}

		if err := query.Find(&list).Error; err != nil {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "no such table") {
				_ = db.AutoMigrate(&models.Station{})
				if err2 := db.Order("id asc").Find(&list).Error; err2 == nil {
					// Format minimal response with available fields
					type stationResponse struct {
						ID     uint   `json:"id"`
						Name   string `json:"name"`
						TypeID uint   `json:"type_id"`
					}

					response := make([]stationResponse, 0, len(list))
					for _, s := range list {
						response = append(response, stationResponse{
							ID:     s.ID,
							Name:   s.Name,
							TypeID: s.StationTypeID,
						})
					}

					return c.JSON(fiber.Map{"data": response})
				}
			}
			return c.JSON(fiber.Map{"data": []models.Station{}})
		}
		return c.JSON(fiber.Map{"data": list})
	}
}

// StationBatch handles batch operations for Station Types
func StationBatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type StationResource struct {
			Name     string `json:"name"`
			Type     string `json:"type"`
			Location string `json:"location"`
		}
		var resources []StationResource
		if err := c.BodyParser(&resources); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}
		// TODO: Insert resources into DB (placeholder logic)
		return c.JSON(fiber.Map{"message": "Batch insert successful", "count": len(resources)})
	}
}
