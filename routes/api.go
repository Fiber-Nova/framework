package routes

import (
	"backend-meta-data/controllers"
	"backend-meta-data/middleware"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterAPIRoutes registers the /api group endpoints
func RegisterAPIRoutes(app *fiber.App, dbConn *sql.DB, gormDB *gorm.DB) {
	api := app.Group("/api")

	// Instrument Types
	api.Post("/instrument-types", middleware.AuthMiddleware, controllers.CreateInstrumentType(gormDB))
	api.Get("/instrument-types", controllers.ListInstrumentTypes(gormDB))
	api.Patch("/instrument-types/:id", middleware.AuthMiddleware, controllers.UpdateInstrumentType(gormDB))

	// Users
	api.Post("/users", controllers.CreateUser(gormDB))
	api.Get("/users", controllers.ListUsers(gormDB))
	api.Patch("/users/:id", controllers.UpdateUser(gormDB))
	api.Patch("/users/profile", middleware.AuthMiddleware, controllers.UpdateUser(gormDB))

	// Stations
	api.Post("/station/batch", controllers.StationBatch())
	api.Get("/station", controllers.ListStations(gormDB))

	// Stores
	api.Get("/stores", controllers.ListStores(gormDB))
	api.Post("/stores", controllers.CreateStore(gormDB))

	// Templates
	api.Get("/templates/maint_notice", controllers.ListMaintNoticeTemplates())
	api.Get("/templates/maint_notice/:name", controllers.GetMaintNoticeTemplate())
	api.Post("/maint-notices", controllers.CreateMaintNotice())
	api.Get("/templates", controllers.ListTemplates())
	api.Get("/stn", middleware.AuthMiddleware, controllers.ListStations(gormDB))
	api.Get("/instruments", controllers.ListInstruments(gormDB))
}
