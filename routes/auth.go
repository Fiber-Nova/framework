package routes

import (
	"backend-meta-data/controllers"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterAuthRoutes registers authentication-related endpoints
func RegisterAuthRoutes(app *fiber.App, dbConn *sql.DB, gormDB *gorm.DB) {
	app.Post("/login", controllers.Login(dbConn))
	app.Post("/logout", controllers.Logout())
	app.Get("/me", controllers.Me(dbConn))
	app.Get("/auth/sso", controllers.SSOController(gormDB))
	app.Post("/auth/ad-login", controllers.ADLoginHandler)
}
