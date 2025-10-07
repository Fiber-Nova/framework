package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes aggregates sub-route registrations
func RegisterRoutes(app *fiber.App, dbConn *sql.DB, gormDB *gorm.DB) {
	RegisterHealthRoutes(app, dbConn)
	RegisterAuthRoutes(app, dbConn, gormDB)
	RegisterAPIRoutes(app, dbConn, gormDB)
	RegisterExportRoutes(app, dbConn)
}
