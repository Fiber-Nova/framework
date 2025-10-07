package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// This file exists to avoid empty package compilation errors.
// Export-related routes are registered in api.go.
func RegisterExportRoutes(app *fiber.App, dbConn *sql.DB) {

}
