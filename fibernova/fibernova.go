// Package fibernova provides a full-stack, Laravel-inspired web framework built on Go Fiber
package fibernova

import (
	"github.com/CasperHK/FiberNova/app"
	"github.com/CasperHK/FiberNova/config"
	"github.com/CasperHK/FiberNova/routing"
	"github.com/gofiber/fiber/v2"
)

// Application is the main FiberNova application
type Application struct {
	*app.Application
	router *routing.Router
}

// New creates a new FiberNova application
func New(cfg *config.Config) *Application {
	appConfig := &app.Config{
		AppName:     cfg.App.Name,
		AppEnv:      cfg.App.Env,
		AppDebug:    cfg.App.Debug,
		AppPort:     cfg.Server.Port,
		AppKey:      cfg.App.Key,
		DatabaseURL: "",
	}
	
	baseApp := app.New(appConfig)
	
	return &Application{
		Application: baseApp,
		router:      routing.New(baseApp.Fiber()),
	}
}

// Router returns the application router
func (app *Application) Router() *routing.Router {
	return app.router
}

// Use adds middleware to the application
func (app *Application) Use(middleware ...fiber.Handler) {
	for _, m := range middleware {
		app.Fiber().Use(m)
	}
}

// Static serves static files from a directory
func (app *Application) Static(prefix, root string) {
	app.Fiber().Static(prefix, root)
}
