package main

import (
	"log"

	"github.com/CasperHK/FiberNova/fibernova"
	"github.com/CasperHK/FiberNova/config"
	"github.com/CasperHK/FiberNova/middleware"
	"github.com/CasperHK/FiberNova/routing"
	"github.com/CasperHK/FiberNova/validation"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Create FiberNova application
	app := fibernova.New(cfg)
	
	// Apply global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	
	// Define routes
	router := app.Router()
	
	// Simple route
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to FiberNova!",
			"version": "1.0.0",
		})
	})
	
	// Route with validation
	router.Post("/users", func(c *fiber.Ctx) error {
		type CreateUserRequest struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		
		// Validate request
		validator := validation.New()
		validator.Required("name", req.Name).Min("name", req.Name, 3)
		validator.Required("email", req.Email).Email("email", req.Email)
		
		if validator.Fails() {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": validator.Errors(),
			})
		}
		
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user": fiber.Map{
				"name":  req.Name,
				"email": req.Email,
			},
		})
	})
	
	// API group with routes
	router.Group("/api", func(group *routing.RouterGroup) {
		group.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "healthy",
			})
		})
		
		group.Get("/version", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"version": "1.0.0",
				"framework": "FiberNova",
			})
		})
	})
	
	// Protected routes group
	router.Group("/admin", func(group *routing.RouterGroup) {
		// Apply auth middleware to this group
		group.Middleware(middleware.Auth(cfg.App.Key))
		
		group.Get("/dashboard", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Welcome to admin dashboard",
			})
		})
	})
	
	// Start the server
	log.Printf("Starting FiberNova on port %s...", cfg.Server.Port)
	if err := app.Listen(); err != nil {
		log.Fatal(err)
	}
}
