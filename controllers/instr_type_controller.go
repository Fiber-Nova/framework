package controllers

import (
	"backend-meta-data/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateInstrumentType handles POST /api/instrument-types
func CreateInstrumentType(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Code        string `json:"code"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}
		if req.Code == "" || req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "code and name are required"})
		}
		it := models.InstrumentType{Code: req.Code, Name: req.Name, Description: req.Description}
		if err := db.Create(&it).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create instrument type"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": it})
	}
}

// ListInstrumentTypes handles GET /api/instrument-types and shapes rows for the UI
func ListInstrumentTypes(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []models.InstrumentType
		if err := db.Order("name asc").Find(&items).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch instrument types"})
		}
		return c.JSON(fiber.Map{"data": items})
	}
}

// UpdateInstrumentType handles PUT /api/instrument-types/:id
func UpdateInstrumentType(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var req struct {
			Code        *string `json:"code"`
			Name        *string `json:"name"`
			Category    *string `json:"category"`
			Status      *string `json:"status"`
			Description *string `json:"description"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}

		var it models.InstrumentType
		if err := db.First(&it, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "instrument type not found"})
		}

		if req.Code != nil {
			it.Code = *req.Code
		}
		if req.Name != nil {
			it.Name = *req.Name
		}
		if req.Category != nil {
			it.Category = *req.Category
		}
		if req.Status != nil {
			it.Status = *req.Status
		}
		if req.Description != nil {
			it.Description = *req.Description
		}

		if err := db.Save(&it).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update instrument type"})
		}

		return c.JSON(fiber.Map{"data": it})
	}
}
