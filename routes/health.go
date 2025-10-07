package routes

import (
	"backend-meta-data/controllers"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// RegisterHealthRoutes registers root, dbcheck and health endpoints
func RegisterHealthRoutes(app *fiber.App, dbConn *sql.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/dbcheck", func(c *fiber.Ctx) error {
		if err := dbConn.Ping(); err != nil {
			return c.Status(500).SendString("DB connection failed")
		}
		return c.SendString("DB connection successful")
	})

	app.Get("/healthz", controllers.Healthz())
}
