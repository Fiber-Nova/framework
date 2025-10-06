package app

import (
	"github.com/gofiber/fiber/v2"
)

// Application is the core of the FiberNova framework
type Application struct {
	fiber     *fiber.App
	config    *Config
	container *ServiceContainer
}

// Config holds application configuration
type Config struct {
	AppName     string
	AppEnv      string
	AppDebug    bool
	AppPort     string
	AppKey      string
	DatabaseURL string
}

// ServiceContainer holds application services
type ServiceContainer struct {
	services map[string]interface{}
}

// New creates a new FiberNova application instance
func New(config *Config) *Application {
	if config == nil {
		config = &Config{
			AppName:  "FiberNova",
			AppEnv:   "development",
			AppDebug: true,
			AppPort:  "3000",
		}
	}

	return &Application{
		fiber:  fiber.New(fiber.Config{
			AppName:      config.AppName,
			ServerHeader: "FiberNova",
		}),
		config: config,
		container: &ServiceContainer{
			services: make(map[string]interface{}),
		},
	}
}

// Fiber returns the underlying Fiber instance
func (app *Application) Fiber() *fiber.App {
	return app.fiber
}

// Config returns the application config
func (app *Application) Config() *Config {
	return app.config
}

// Container returns the service container
func (app *Application) Container() *ServiceContainer {
	return app.container
}

// Listen starts the HTTP server
func (app *Application) Listen() error {
	return app.fiber.Listen(":" + app.config.AppPort)
}

// Shutdown gracefully shuts down the server
func (app *Application) Shutdown() error {
	return app.fiber.Shutdown()
}

// Bind registers a service in the container
func (c *ServiceContainer) Bind(name string, service interface{}) {
	c.services[name] = service
}

// Get retrieves a service from the container
func (c *ServiceContainer) Get(name string) (interface{}, bool) {
	service, exists := c.services[name]
	return service, exists
}
