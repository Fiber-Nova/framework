package controllers

import (
	"backend-meta-data/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ListStores handles GET /api/stores
func ListStores(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		search := c.Query("q")
		page := c.QueryInt("page", 1)
		size := c.QueryInt("pageSize", 20)
		sortField := c.Query("sortField", "name")
		sortOrder := c.Query("sortOrder", "asc")

		rows, total, err := models.ListInventoryStores(db, search, page, size, sortField, sortOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch stores"})
		}
		return c.JSON(fiber.Map{"data": rows, "total": total, "page": page, "pageSize": size})
	}
}

// CreateStore handles POST /api/stores
func CreateStore(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.InventoryStore
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}
		if err := models.CreateInventoryStore(db, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": req})
	}
}
